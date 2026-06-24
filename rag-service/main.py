import os
import logging
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import StreamingResponse
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

from models import IngestRequest, IngestResponse, QueryRequest
from ingestion import ingest_document
from retrieval import retrieve_and_generate

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

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
