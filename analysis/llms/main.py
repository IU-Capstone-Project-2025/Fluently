import os
import asyncio
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from starlette_prometheus import PrometheusMiddleware, metrics
from pydantic import BaseModel
from typing import List, Optional
import uvicorn

from api import AIService

# Global AI service instance
ai_service: AIService = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialize AI service at startup"""
    global ai_service
    print("Initializing AI service...")
    try:
        ai_service = AIService()
        await ai_service.initialize()
        print("AI service initialized successfully!")
    except Exception as e:
        print(f"Failed to initialize AI service: {e}")
        print("The API will still start, but AI functionality will be limited")
        ai_service = None
    
    yield
    
    # Cleanup if needed
    if ai_service:
        del ai_service

app = FastAPI(
    title="Fluently LLM API",
    description="AI-powered conversation service for language learning",
    version="1.0.0",
    lifespan=lifespan
)

# Add Prometheus middleware
app.add_middleware(PrometheusMiddleware)
app.add_route("/metrics", metrics)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure as needed for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Request/Response models
class Message(BaseModel):
    role: str  # "user", "assistant", or "system"
    content: str

class ChatRequest(BaseModel):
    model_config = {"protected_namespaces": ()}
    
    messages: List[Message]
    model_type: Optional[str] = "balanced"  # "fast" or "balanced"
    max_tokens: Optional[int] = None
    temperature: Optional[float] = None

class ChatResponse(BaseModel):
    model_config = {"protected_namespaces": ()}
    
    response: str
    model_used: Optional[str] = None

class HealthResponse(BaseModel):
    status: str
    ai_service_ready: bool
    available_providers: Optional[List[str]] = None
    message: Optional[str] = None

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    if ai_service is None:
        return HealthResponse(
            status="degraded",
            ai_service_ready=False,
            available_providers=[],
            message="AI service not initialized"
        )
    
    # Check which providers have valid keys
    available_providers = [
        provider for provider in ai_service.providers 
        if ai_service.providers[provider]["keys"]
    ]
    
    if not available_providers:
        return HealthResponse(
            status="degraded",
            ai_service_ready=False,
            available_providers=available_providers,
            message="No valid API keys configured. Please check GROQ_API_KEYS and GEMINI_API_KEYS environment variables."
        )
    
    return HealthResponse(
        status="healthy",
        ai_service_ready=True,
        available_providers=available_providers,
        message=f"AI service ready with {len(available_providers)} provider(s): {', '.join(available_providers)}"
    )

class ConfigResponse(BaseModel):
    providers: dict
    setup_urls: dict = {
        "groq": "https://console.groq.com/keys",
        "gemini": "https://makersuite.google.com/app/apikey"
    }

@app.get("/config", response_model=ConfigResponse)
async def get_config():
    """Get configuration status"""
    if ai_service is None:
        return ConfigResponse(
            providers={
                "groq": {"configured": False, "keys_count": 0},
                "gemini": {"configured": False, "keys_count": 0}
            }
        )
    
    providers_info = {}
    for provider, provider_data in ai_service.providers.items():
        keys = provider_data["keys"]
        providers_info[provider] = {
            "configured": len(keys) > 0,
            "keys_count": len(keys),
            "has_valid_keys": len([k for k in keys if k not in ["placeholder", "your_key_here", ""]]) > 0
        }
    
    return ConfigResponse(providers=providers_info)

@app.post("/chat", response_model=ChatResponse)
async def chat_completion(request: ChatRequest):
    """
    Generate AI response for conversation
    """
    if ai_service is None:
        raise HTTPException(status_code=503, detail="AI service not initialized")
    
    try:
        # Convert Pydantic models to dict format expected by AI service
        messages = [{"role": msg.role, "content": msg.content} for msg in request.messages]
        
        # Prepare kwargs for AI service
        kwargs = {}
        if request.max_tokens:
            kwargs["max_tokens"] = request.max_tokens
        if request.temperature:
            kwargs["temperature"] = request.temperature
        
        # Get AI response
        response = await ai_service.chat_completion(
            messages=messages,
            model_type=request.model_type,
            **kwargs
        )
        
        return ChatResponse(response=response)
        
    except Exception as e:
        error_msg = str(e)
        if "No providers are available" in error_msg:
            raise HTTPException(
                status_code=503, 
                detail="AI service unavailable: No valid API keys configured. Please check GROQ_API_KEYS and GEMINI_API_KEYS environment variables."
            )
        elif "All available providers failed" in error_msg:
            raise HTTPException(
                status_code=503, 
                detail="AI service unavailable: All configured providers failed. Please check your API keys are valid and have sufficient quota."
            )
        else:
            raise HTTPException(status_code=500, detail=f"AI service error: {error_msg}")

@app.post("/chat/simple")
async def simple_chat(message: str, model_type: str = "balanced"):
    """
    Simple chat endpoint for quick single-message conversations
    """
    if ai_service is None:
        raise HTTPException(status_code=503, detail="AI service not initialized")
    
    try:
        messages = [{"role": "user", "content": message}]
        response = await ai_service.chat_completion(
            messages=messages,
            model_type=model_type
        )
        return {"response": response}
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"AI service error: {str(e)}")

@app.get("/status")
async def get_status():
    """Get detailed service status"""
    if ai_service is None:
        return {"status": "Service not initialized"}
    
    return {
        "status": "ready",
        "providers": {
            "groq": {
                "keys_available": len(ai_service.providers["groq"]["keys"]),
                "models": ai_service.providers["groq"]["models"]
            },
            "gemini": {
                "keys_available": len(ai_service.providers["gemini"]["keys"]),
                "models": ai_service.providers["gemini"]["models"]
            }
        }
    }

def main():
    """Main entry point for development"""
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8003,
        reload=True,
        log_level="info"
    )

if __name__ == "__main__":
    main()
