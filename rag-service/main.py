import os
import logging
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import StreamingResponse
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

from models import IngestRequest, IngestResponse, QueryRequest, AnalysisRequest, MealPlanRequest, WorkoutPlanRequest
from ingestion import ingest_document, get_supabase_client, download_file_from_supabase, parse_pdf
from retrieval import retrieve_and_generate
from analyzer import analyze_document_text
from meal_planner import generate_meal_plan
from workout_planner import generate_workout_plan


# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
)
logger = logging.getLogger(__name__)

app = FastAPI(title="FitMind RAG Microservice", version="1.0.0")

@app.get("/health")
def health_check():
    return {"status": "healthy"}

@app.post("/ingest", response_model=IngestResponse)
def ingest_endpoint(request: IngestRequest):
    try:
        logger.info(f"Received ingestion request for doc_id={request.doc_id}")
        num_chunks = ingest_document(
            user_id=request.user_id,
            doc_id=request.doc_id,
            storage_path=request.storage_path
        )
        return IngestResponse(
            success=True,
            message=f"Document successfully ingested into {num_chunks} chunks."
        )
    except Exception as e:
        logger.exception("Failed to ingest document")
        raise HTTPException(
            status_code=500,
            detail=f"Ingestion failed: {str(e)}"
        )

@app.post("/query")
def query_endpoint(request: QueryRequest):
    logger.info(f"Received query request for user_id={request.user_id}, doc_id={request.doc_id}")
    try:
        # Stream the results via SSE
        generator = retrieve_and_generate(
            query=request.query,
            doc_id=request.doc_id,
            user_id=request.user_id,
            chat_history=request.chat_history
        )
        return StreamingResponse(
            generator,
            media_type="text/event-stream",
            headers={
                "Cache-Control": "no-cache",
                "Connection": "keep-alive",
                "Content-Type": "text/event-stream"
            }
        )
    except Exception as e:
        logger.exception("Failed to query document")
        raise HTTPException(
            status_code=500,
            detail=f"Query failed: {str(e)}"
        )

@app.post("/analyze")
def analyze_endpoint(request: AnalysisRequest):
    logger.info(f"Received analysis request for doc_id={request.doc_id}")
    try:
        supabase = get_supabase_client()
        
        # 1. Fetch document metadata to get the storage path
        doc_res = supabase.table("documents").select("storage_path").eq("id", request.doc_id).execute()
        if not doc_res.data:
            raise HTTPException(status_code=404, detail="Document not found")
        storage_path = doc_res.data[0]["storage_path"]
        
        # 2. Download and parse PDF
        file_bytes = download_file_from_supabase(storage_path)
        documents_list = parse_pdf(file_bytes, request.user_id, request.doc_id)
        full_text = "\n".join([doc.text for doc in documents_list])
        
        # 3. Analyze text
        analysis_result = analyze_document_text(full_text)
        return {"success": True, "data": analysis_result}
    except HTTPException:
        raise
    except Exception as e:
        logger.exception("Failed to analyze document")
        raise HTTPException(
            status_code=500,
            detail=f"Analysis failed: {str(e)}"
        )

@app.post("/meal-plan")
def meal_plan_endpoint(request: MealPlanRequest):
    logger.info(f"Received meal-plan request for user_id={request.user_id}, doc_id={request.doc_id}")
    try:
        supabase = get_supabase_client()
        
        # 1. Fetch analysis
        analysis_res = supabase.table("document_analyses").select("*").eq("doc_id", request.doc_id).execute()
        analysis = analysis_res.data[0] if analysis_res.data else {"metrics": [], "summary": "No document analysis available."}
        
        # 2. Fetch user goals
        goals_res = supabase.table("goals").select("*").eq("user_id", request.user_id).execute()
        goals = goals_res.data if goals_res.data else []
        
        # 3. Fetch latest BMI metric
        bmi_res = supabase.table("health_metrics").select("bmi").eq("user_id", request.user_id).order("recorded_at", desc=True).limit(1).execute()
        bmi = bmi_res.data[0]["bmi"] if bmi_res.data else None
        
        # 4. Generate meal plan
        plan = generate_meal_plan(analysis, goals, bmi, request.dietary_preference)
        
        # 5. Store meal plan in database
        plan_data = {
            "user_id": request.user_id,
            "doc_id": request.doc_id,
            "reasoning": plan.get("reasoning", ""),
            "daily_calories_target": plan.get("daily_calories_target", 2000),
            "protein_target_g": plan.get("protein_target_g", 120),
            "days": plan.get("days", []),
            "weekly_notes": plan.get("weekly_notes", ""),
            "dietary_preference": request.dietary_preference
        }
        
        logger.info(f"Saving generated meal plan to Supabase for user_id={request.user_id}")
        db_res = supabase.table("meal_plans").insert(plan_data).execute()
        if not db_res.data:
            raise Exception("Failed to insert meal plan record")
            
        return {"success": True, "data": db_res.data[0]}
    except Exception as e:
        logger.exception("Failed to generate meal plan")
        raise HTTPException(
            status_code=500,
            detail=f"Meal plan generation failed: {str(e)}"
        )

@app.post("/workout-plan")
def workout_plan_endpoint(request: WorkoutPlanRequest):
    logger.info(f"Received workout-plan request for user_id={request.user_id}")
    try:
        supabase = get_supabase_client()
        
        # 1. Fetch user goals
        goals_res = supabase.table("goals").select("*").eq("user_id", request.user_id).execute()
        goals = goals_res.data if goals_res.data else []
        
        # 2. Fetch latest BMI metric
        bmi_res = supabase.table("health_metrics").select("bmi").eq("user_id", request.user_id).order("recorded_at", desc=True).limit(1).execute()
        bmi = bmi_res.data[0]["bmi"] if bmi_res.data else None
        
        # 3. Generate workout plan
        plan = generate_workout_plan(bmi, goals, request.fitness_level, request.equipment, request.days_per_week)
        
        # 4. Store workout plan in database
        plan_data = {
            "user_id": request.user_id,
            "program_name": plan.get("program_name", "FitMind Program"),
            "duration_weeks": plan.get("duration_weeks", 4),
            "reasoning": plan.get("reasoning", ""),
            "weekly_schedule": plan.get("weekly_schedule", []),
            "progression_notes": plan.get("progression_notes", ""),
            "safety_notes": plan.get("safety_notes", ""),
            "fitness_level": request.fitness_level,
            "equipment": request.equipment,
            "days_per_week": request.days_per_week
        }
        
        logger.info(f"Saving generated workout plan to Supabase for user_id={request.user_id}")
        db_res = supabase.table("workout_plans").insert(plan_data).execute()
        if not db_res.data:
            raise Exception("Failed to insert workout plan record")
            
        return {"success": True, "data": db_res.data[0]}
    except Exception as e:
        logger.exception("Failed to generate workout plan")
        raise HTTPException(
            status_code=500,
            detail=f"Workout plan generation failed: {str(e)}"
        )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
