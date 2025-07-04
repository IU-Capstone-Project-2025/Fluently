from pydantic import BaseModel, Field
from typing import List


class DistractorRequest(BaseModel):
    """Request schema for distractor generation"""
    sentence: str = Field(..., min_length=1, max_length=1000, description="The sentence containing the target word")
    word: str = Field(..., min_length=1, max_length=100, description="The target word to generate distractors for")
    
    class Config:
        json_schema_extra = {
            "example": {
                "sentence": "The cat caught the mouse in the kitchen",
                "word": "mouse"
            }
        }


class DistractorResponse(BaseModel):
    """Response schema for distractor generation"""
    pick_options: List[str] = Field(..., description="List of distractor words including the correct answer")
    
    class Config:
        json_schema_extra = {
            "example": {
                "pick_options": ["mouse", "rat", "bird", "fish"]
            }
        }


class ErrorResponse(BaseModel):
    """Error response schema"""
    detail: str = Field(..., description="Error message")
    error_code: str = Field(None, description="Error code for client handling")


class TTSRequest(BaseModel):
    """Request schema for text-to-speech generation"""
    text: str = Field(..., min_length=1, max_length=200, description="The text to convert to speech")
    
    class Config:
        json_schema_extra = {
            "example": {
                "text": "Hello world"
            }
        } 