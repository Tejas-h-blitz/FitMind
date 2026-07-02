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

def generate_workout_plan(bmi: float, goals: list, fitness_level: str, equipment: str, days_per_week: int) -> dict:
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        logger.error("OPENAI_API_KEY is not set in environment variables")
        raise ValueError("OPENAI_API_KEY must be set in the environment variables")

    client = OpenAI(
        api_key=api_key,
        base_url="https://api.groq.com/openai/v1"
    )

    # 1. Format goals
    goals_list = [g.get("title") for g in goals if g.get("status") == "active"]
    goals_str = ", ".join(goals_list) if goals_list else "General fitness"

    # 2. Format BMI category
    bmi_category = get_bmi_category(bmi) if bmi else "Unknown"
    bmi_val = f"{bmi:.1f}" if bmi else "Not provided"

    # 3. Build prompt exactly as requested
    prompt = (
        "You are a certified personal trainer AI. Generate a personalized workout\n"
        "plan based on the following:\n"
        f"FITNESS LEVEL: {fitness_level} (beginner/intermediate/advanced)\n"
        f"EQUIPMENT: {equipment} (none/home/full_gym)\n"
        f"DAYS PER WEEK: {days_per_week}\n"
        f"BMI: {bmi_val} (Category: {bmi_category})\n"
        f"HEALTH GOALS: {goals_str}\n"
        "Return ONLY valid JSON. No markdown. No preamble.\n"
        "Schema:\n"
        "{\n"
        '  "program_name": "catchy program name",\n'
        '  "duration_weeks": 4,\n'
        '  "reasoning": "why this program suits this user",\n'
        '  "weekly_schedule": [\n'
        "    {\n"
        '      "day": "Monday",\n'
        '      "focus": "e.g. Upper Body Strength",\n'
        '      "is_rest_day": false,\n'
        '      "exercises": [\n'
        "        {\n"
        '          "name": "exercise name",\n'
        '          "sets": number,\n'
        '          "reps": "e.g. 8-12 or 30 seconds",\n'
        '          "rest_seconds": number,\n'
        '          "instructions": "brief form cue",\n'
        '          "muscle_groups": ["chest", "triceps"],\n'
        '          "modification_easier": "easier alternative",\n'
        '          "modification_harder": "harder alternative"\n'
        "        }\n"
        "      ],\n"
        '      "estimated_duration_minutes": number,\n'
        '      "warmup": "brief warmup description",\n'
        '      "cooldown": "brief cooldown description"\n'
        "    }\n"
        "  ],\n"
        '  "progression_notes": "how to progress after 4 weeks",\n'
        '  "safety_notes": "important safety reminders"\n'
        "}\n"
    )

    try:
        logger.info("Sending request to Groq for workout planning")
        response = client.chat.completions.create(
            model="llama3-70b-8192",
            messages=[{"role": "user", "content": prompt}],
            response_format={"type": "json_object"},
            temperature=0.2
        )
        
        raw_content = response.choices[0].message.content
        logger.info("Successfully received workout plan from GPT-4o")
        return json.loads(raw_content)

    except Exception as e:
        logger.exception("Error during workout planning via OpenAI GPT-4o")
        raise e
