#!/bin/bash
# Quick test to check if our changes resolved the Groq initialization issue

echo "🧪 Quick test for Groq initialization issue..."

# Create a temporary test directory
TEST_DIR="/tmp/llm_test"
mkdir -p $TEST_DIR

# Copy files to test directory
cp /home/timofey/Desktop/Projects/Fluently/Fluently-fork/analysis/llms/requirements.txt $TEST_DIR/
cp /home/timofey/Desktop/Projects/Fluently/Fluently-fork/analysis/llms/api.py $TEST_DIR/
cp /home/timofey/Desktop/Projects/Fluently/Fluently-fork/analysis/llms/main.py $TEST_DIR/

# Create a minimal test script
cat > $TEST_DIR/minimal_test.py << 'EOF'
#!/usr/bin/env python3
import sys
import os

# Add dummy API keys
os.environ["GROQ_API_KEYS"] = "dummy_key"
os.environ["GEMINI_API_KEYS"] = "dummy_key"

try:
    from groq import Groq
    print("✅ Groq import successful")
    
    # Test client creation
    try:
        client = Groq(api_key="dummy_key")
        print("✅ Groq client creation successful")
    except Exception as e:
        print(f"❌ Groq client creation failed: {e}")
        sys.exit(1)
        
    # Test API service
    try:
        from api import AIService
        print("✅ AIService import successful")
        
        import asyncio
        
        async def test_init():
            service = AIService()
            await service.initialize()
            print("✅ AIService initialization successful")
            
        asyncio.run(test_init())
        
    except Exception as e:
        print(f"❌ AIService test failed: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
        
    print("🎉 All tests passed!")
    
except ImportError as e:
    print(f"❌ Import failed: {e}")
    sys.exit(1)
EOF

# Create a minimal Dockerfile for testing
cat > $TEST_DIR/Dockerfile << 'EOF'
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "minimal_test.py"]
EOF

echo "📦 Building test Docker image..."
cd $TEST_DIR
docker build -t groq-test:latest .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "🏃 Running test container..."
docker run --rm groq-test:latest

TEST_RESULT=$?

# Cleanup
rm -rf $TEST_DIR

if [ $TEST_RESULT -eq 0 ]; then
    echo "✅ Test passed! The Groq initialization issue should be resolved."
else
    echo "❌ Test failed. The issue persists."
    exit 1
fi
