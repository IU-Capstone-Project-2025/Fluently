import asyncio
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

from models.schemas import DistractorRequest, DistractorResponse
from services.distractor_service import DistractorService
from api.routes import router


# Global service instance for optimal latency
distractor_service: DistractorService = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialize model at startup for optimal latency"""
    global distractor_service
    print("Loading BERT model...")
    distractor_service = DistractorService()
    await distractor_service.initialize()
    print("Model loaded successfully!")
    yield
    # Cleanup if needed
    del distractor_service


app = FastAPI(
    title="Distractor Generator API",
    description="High-performance API for generating word distractors using BERT",
    version="1.0.0",
    lifespan=lifespan
)

# Add CORS middleware for frontend integration
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure as needed for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


def get_distractor_service() -> DistractorService:
    """Dependency injection for distractor service"""
    if distractor_service is None:
        raise HTTPException(status_code=503, detail="Service not initialized")
    return distractor_service


# Override the dependency in routes
from api.routes import get_distractor_service as routes_get_service
app.dependency_overrides[routes_get_service] = get_distractor_service

# Include API routes
app.include_router(router)


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "model_loaded": distractor_service is not None}


def main():
    """Main entry point for Poetry script"""
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8001,
        reload=False,  # Disable reload for better performance
        workers=1,  # Single worker to avoid model loading overhead
    )


if __name__ == "__main__":
    main()
