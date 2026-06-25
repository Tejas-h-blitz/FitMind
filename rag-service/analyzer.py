import os
import json
import logging
from openai import OpenAI

logger = logging.getLogger(__name__)

def analyze_document_text(text: str) -> dict:
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        logger.error("OPENAI_API_KEY is not set in environment variables")
        raise ValueError("OPENAI_API_KEY must be set in the environment variables")

    client = OpenAI(api_key=api_key)

    # Building the exact prompt requested
    prompt = (
        "You are a medical document analyzer. Extract health metrics from the\n"
        "following document text. Return ONLY valid JSON, no markdown, no\n"
        "preamble. Use null for any metric not found in the document.\n"
        "Schema:\n"
        "{\n"
        '  "metrics": [\n'
        "    {\n"
        '      "name": "metric name",\n'
        '      "value": numeric_value,\n'
        '      "unit": "unit string",\n'
        '      "status": "normal|low|high|borderline",\n'
        '      "reference_range": "e.g. 13.5-17.5",\n'
        '      "plain_english": "one sentence explanation"\n'
        "    }\n"
        "  ],\n"
        '  "summary": "2-3 sentence plain English overall summary",\n'
        '  "overall_status": "good|needs_attention|concerning"\n'
        "}\n\n"
        "Metrics to extract if present:\n"
        "- hemoglobin\n"
        "- blood_sugar_fasting\n"
        "- blood_sugar_postprandial\n"
        "- total_cholesterol\n"
        "- hdl_cholesterol\n"
        "- ldl_cholesterol\n"
        "- triglycerides\n"
        "- vitamin_d\n"
        "- vitamin_b12\n"
        "- tsh\n"
        "- hemoglobin_a1c\n"
        "- systolic_bp\n"
        "- diastolic_bp\n"
        "- creatinine\n"
        "- uric_acid\n\n"
        f"Document text to analyze:\n{text}"
    )

    try:
        logger.info("Sending request to GPT-4o for document analysis")
        response = client.chat.completions.create(
            model="gpt-4o",
            messages=[{"role": "user", "content": prompt}],
            response_format={"type": "json_object"},
            temperature=0.1
        )
        
        raw_content = response.choices[0].message.content
        logger.info("Successfully received analysis from GPT-4o")
        
        analysis_result = json.loads(raw_content)
        
        # Post-process: Filter out metrics where value is null/None, or handle appropriately.
        # The prompt says: "Use null for any metric not found in the document."
        # If the metric is not found, the LLM might return it with value = null.
        # We should preserve or clean them up so the frontend gets a clean array of found metrics.
        if "metrics" in analysis_result and isinstance(analysis_result["metrics"], list):
            filtered_metrics = []
            for item in analysis_result["metrics"]:
                if item and item.get("value") is not None and item.get("name") is not None:
                    filtered_metrics.append(item)
            analysis_result["metrics"] = filtered_metrics

        return analysis_result

    except Exception as e:
        logger.exception("Error during document analysis via OpenAI GPT-4o")
        raise e
