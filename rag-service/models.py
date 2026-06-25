from pydantic import BaseModel, Field
from typing import List, Dict, Any, Optional

class IngestRequest(BaseModel):
    doc_id: str = Field(..., description="The unique UUID of the document")
    user_id: str = Field(..., description="The user's UUID")
    storage_path: str = Field(..., description="The Supabase Storage path of the PDF")

class IngestResponse(BaseModel):
    success: bool
    message: str

class MessageParam(BaseModel):
    role: str = Field(..., description="user or assistant")
    content: str = Field(..., description="Message text")

class QueryRequest(BaseModel):
    query: str = Field(..., description="The user's question")
    doc_id: str = Field(..., description="The document UUID")
    user_id: str = Field(..., description="The user's UUID")
    chat_history: List[MessageParam] = Field(default=[], description="List of previous messages")

class SourceChunk(BaseModel):
    text: str
    page_number: int

class AnalysisRequest(BaseModel):
    doc_id: str
    user_id: str

class MealPlanRequest(BaseModel):
    user_id: str
    doc_id: str
    dietary_preference: str  # vegetarian, non-vegetarian, vegan

class WorkoutPlanRequest(BaseModel):
    user_id: str
    fitness_level: str  # beginner, intermediate, advanced
    equipment: str  # none, home, full_gym
    days_per_week: int

