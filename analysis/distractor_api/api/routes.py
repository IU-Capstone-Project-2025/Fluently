from fastapi import APIRouter, HTTPException, Depends, status
from fastapi.responses import JSONResponse
import time
import logging

from models.schemas import DistractorRequest, DistractorResponse, ErrorResponse
from services.distractor_service import DistractorService

logger = logging.getLogger(__name__)

router = APIRouter(
    prefix="/api/v1",
    tags=["distractors"]
)


def get_distractor_service() -> DistractorService:
    """Dependency injection placeholder - will be overridden in main.py"""
    pass


@router.post(
    "/generate-distractors",
    response_model=DistractorResponse,
    responses={
        400: {"model": ErrorResponse, "description": "Bad Request"},
        422: {"model": ErrorResponse, "description": "Validation Error"},
        503: {"model": ErrorResponse, "description": "Service Unavailable"},
    },
    summary="Generate word distractors",
    description="Generate distractor words for a given target word in a sentence using BERT model"
)
async def generate_distractors(
    request: DistractorRequest,
    service: DistractorService = Depends(get_distractor_service)
) -> DistractorResponse:
    """
    Generate distractors for a target word in a sentence.
    
    The API returns the original word along with generated distractors,
    all shuffled randomly to create multiple choice options.
    """
    start_time = time.time()
    
    try:
        # Input validation
        if not request.sentence.strip():
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Sentence cannot be empty"
            )
        
        if not request.word.strip():
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Target word cannot be empty"
            )
        
        # Check if word exists in sentence (case insensitive)
        if request.word.lower() not in request.sentence.lower():
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Target word '{request.word}' not found in the sentence"
            )
        
        # Generate distractors
        pick_options = await service.generate_distractors(
            sentence=request.sentence.strip(),
            target_word=request.word.strip(),
            num_distractors=3  # Fixed to 3 distractors + 1 correct = 4 total
        )
        
        # Ensure we have at least the original word
        if not pick_options:
            pick_options = [request.word.strip()]
        
        processing_time = time.time() - start_time
        logger.info(f"Generated distractors in {processing_time:.3f}s for word: {request.word}")
        
        return DistractorResponse(pick_options=pick_options)
        
    except HTTPException:
        # Re-raise HTTP exceptions
        raise
    except Exception as e:
        logger.error(f"Unexpected error generating distractors: {e}")
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail="Internal service error occurred"
        )


@router.get(
    "/health",
    summary="Service health check",
    description="Check if the distractor generation service is healthy and ready"
)
async def health_check(service: DistractorService = Depends(get_distractor_service)):
    """Health check endpoint"""
    try:
        is_ready = service.initialized if service else False
        
        return {
            "status": "healthy" if is_ready else "initializing",
            "model_loaded": is_ready,
            "timestamp": time.time()
        }
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        return JSONResponse(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            content={"status": "unhealthy", "error": str(e)}
        ) 