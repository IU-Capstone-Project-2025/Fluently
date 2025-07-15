#!/usr/bin/env python3
"""
Test script for the LLM API service
"""

import asyncio
import sys
import os
from api import AIService

async def test_ai_service():
    """Test the AI service initialization and basic functionality"""
    print("Testing AI Service initialization...")
    
    try:
        ai_service = AIService()
        await ai_service.initialize()
        print("✅ AI service initialized successfully")
        
        # Test basic chat completion
        messages = [
            {"role": "user", "content": "Hello, how are you?"}
        ]
        
        print("Testing chat completion...")
        response = await ai_service.chat_completion(messages, model_type="fast")
        print(f"✅ Chat completion successful: {response[:100]}...")
        
        return True
        
    except Exception as e:
        print(f"❌ Error during testing: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    # Test with environment variables
    if not os.getenv("GROQ_API_KEYS") and not os.getenv("GEMINI_API_KEYS"):
        print("⚠️  No API keys found in environment variables")
        print("Set GROQ_API_KEYS or GEMINI_API_KEYS to test the service")
        sys.exit(1)
    
    success = asyncio.run(test_ai_service())
    sys.exit(0 if success else 1)
