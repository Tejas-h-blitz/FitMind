import os
import json
import logging
from openai import OpenAI

logger = logging.getLogger(__name__)

def get_bmi_category(bmi: float) -> str:
    if bmi < 18.5:
        return "Underweight"
    elif bmi < 25.0:
        return "Normal weight"
    elif bmi < 30.0:
        return "Overweight"
    else:
        return "Obese"

def generate_meal_plan(analysis: dict, goals: list, bmi: float, dietary_preference: str) -> dict:
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        logger.error("OPENAI_API_KEY is not set in environment variables")
        raise ValueError("OPENAI_API_KEY must be set in the environment variables")

    client = OpenAI(
        api_key=api_key,
        base_url="https://api.groq.com/openai/v1"
    )

    # 1. Format deficiencies
    metrics = analysis.get("metrics", [])
    deficiencies_list = []
    for m in metrics:
        if m.get("status") in ["low", "high", "borderline"]:
            deficiencies_list.append(f"{m.get('name')} ({m.get('value')} {m.get('unit')} - Status: {m.get('status')})")
    
    deficiencies_str = ", ".join(deficiencies_list) if deficiencies_list else "None detected"
    
    # 2. Format goals
    goals_list = [g.get("title") for g in goals if g.get("status") == "active"]
    goals_str = ", ".join(goals_list) if goals_list else "General health maintenance"

    # 3. Format BMI category
    bmi_category = get_bmi_category(bmi) if bmi else "Unknown"
    bmi_val = f"{bmi:.1f}" if bmi else "Not provided"

    # 4. Extract analysis summary
    analysis_summary = analysis.get("summary", "No prior document analysis available.")

    # 5. Build prompt exactly as requested
    prompt = (
        "You are a certified nutritionist AI. Generate a personalized 7-day meal\n"
        "plan based on the following patient data:\n"
        f"HEALTH ANALYSIS: {analysis_summary}\n"
        f"DEFICIENCIES DETECTED: {deficiencies_str}\n"
        f"HEALTH GOALS: {goals_str}\n"
        f"BMI: {bmi_val} (Category: {bmi_category})\n"
        f"DIETARY PREFERENCE: {dietary_preference}\n"
        "Return ONLY valid JSON. No markdown. No preamble.\n"
        "Schema:\n"
        "{\n"
        '  "reasoning": "2-3 sentences explaining why this plan suits this patient",\n'
        '  "daily_calories_target": number,\n'
        '  "protein_target_g": number,\n'
        '  "days": [\n'
        "    {\n"
        '      "day": "Monday",\n'
        '      "meals": {\n'
        '        "breakfast": {\n'
        '          "name": "meal name",\n'
        '          "ingredients": ["ingredient 1", "ingredient 2"],\n'
        '          "calories": number,\n'
        '          "benefits": "one sentence why this helps their specific condition"\n'
        "        },\n"
        '        "lunch": {\n'
        '          "name": "meal name",\n'
        '          "ingredients": ["ingredient 1", "ingredient 2"],\n'
        '          "calories": number,\n'
        '          "benefits": "one sentence why this helps their specific condition"\n'
        "        },\n"
        '        "dinner": {\n'
        '          "name": "meal name",\n'
        '          "ingredients": ["ingredient 1", "ingredient 2"],\n'
        '          "calories": number,\n'
        '          "benefits": "one sentence why this helps their specific condition"\n'
        "        },\n"
        '        "snacks": {\n'
        '          "name": "meal name",\n'
        '          "ingredients": ["ingredient 1", "ingredient 2"],\n'
        '          "calories": number,\n'
        '          "benefits": "one sentence why this helps their specific condition"\n'
        "        }\n"
        "      },\n"
        '      "daily_tip": "one actionable health tip for the day"\n'
        "    }\n"
        "  ],\n"
        '  "weekly_notes": "overall guidance paragraph"\n'
        "}\n"
    )

    try:
        logger.info("Sending request to Groq for meal planning")
        response = client.chat.completions.create(
            model="llama3-70b-8192",
            messages=[{"role": "user", "content": prompt}],
            response_format={"type": "json_object"},
            temperature=0.2
        )
        
        raw_content = response.choices[0].message.content
        logger.info("Successfully received meal plan from GPT-4o")
        return json.loads(raw_content)

    except Exception as e:
        logger.exception("Error during meal planning via OpenAI GPT-4o")
        raise e
