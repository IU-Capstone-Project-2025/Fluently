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
    
    # Test with dummy API keys if none are provided
    if not os.getenv("GROQ_API_KEYS") and not os.getenv("GEMINI_API_KEYS"):
        print("⚠️  No API keys found in environment variables")
        print("Testing initialization with dummy keys...")
        os.environ["GROQ_API_KEYS"] = "dummy_groq_key"
        os.environ["GEMINI_API_KEYS"] = "dummy_gemini_key"
    
    try:
        ai_service = AIService()
        await ai_service.initialize()
        print("✅ AI service initialized (note: may have failed provider initialization)")
        
        # Check which providers are available
        available_providers = [provider for provider in ai_service.providers if ai_service.providers[provider]["keys"]]
        print(f"📋 Available providers: {available_providers}")
        
        if not available_providers:
            print("❌ No providers available - service initialized but won't work")
            return False
        
        # Test basic chat completion only if we have real API keys
        if os.getenv("GROQ_API_KEYS") and "dummy" not in os.getenv("GROQ_API_KEYS"):
            print("Testing chat completion with Groq...")
            messages = [
                {"role": "user", "content": "Hello, respond with just 'Hi there!'"}
            ]
            
            try:
                response = await ai_service.chat_completion(messages, model_type="fast")
                print(f"✅ Chat completion successful: {response[:100]}...")
                return True
            except Exception as e:
                print(f"❌ Chat completion failed: {e}")
                return False
        else:
            print("⚠️  Skipping chat completion test (no real API keys)")
            return True
        
    except Exception as e:
        print(f"❌ Error during testing: {e}")
        import traceback
        traceback.print_exc()
        return False

async def test_groq_client():
    """Test Groq client initialization specifically"""
    print("\n🔍 Testing Groq client initialization...")
    
    try:
        from groq import Groq
        
        # Try with dummy key
        test_key = "dummy_key_for_testing"
        
        try:
            client = Groq(api_key=test_key)
            print("✅ Groq client created successfully")
            return True
        except Exception as e:
            print(f"❌ Groq client creation failed: {e}")
            print(f"Error type: {type(e).__name__}")
            return False
            
    except ImportError as e:
        print(f"❌ Failed to import Groq: {e}")
        return False

if __name__ == "__main__":
    print("🧪 LLM API Service Test Suite")
    print("=" * 50)
    
    # Test Groq client first
    groq_success = asyncio.run(test_groq_client())
    
    # Test AI service
    service_success = asyncio.run(test_ai_service())
    
    print("\n📊 Test Results:")
    print(f"Groq client: {'✅ PASS' if groq_success else '❌ FAIL'}")
    print(f"AI service: {'✅ PASS' if service_success else '❌ FAIL'}")
    
    if not groq_success:
        print("\n💡 Groq client test failed - consider updating the groq library version")
        print("Try: pip install groq==0.15.0 or pip install groq --upgrade")
    
    overall_success = groq_success and service_success
    sys.exit(0 if overall_success else 1)
