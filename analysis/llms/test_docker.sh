#!/bin/bash
# Test script for LLM API Docker container

echo "🧪 Testing LLM API Docker container..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t fluently-llm-api:test .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "✅ Docker image built successfully"

# Test with dummy environment variables
echo "🔧 Testing with dummy environment variables..."
docker run --rm \
    -e GROQ_API_KEYS="dummy_key" \
    -e GEMINI_API_KEYS="dummy_key" \
    --name fluently-llm-test \
    fluently-llm-api:test &

CONTAINER_PID=$!

# Wait for container to start
echo "⏳ Waiting for container to start..."
sleep 10

# Check if container is still running
if ! docker ps | grep -q fluently-llm-test; then
    echo "❌ Container failed to start properly"
    docker logs fluently-llm-test 2>/dev/null || echo "No logs available"
    exit 1
fi

# Test health endpoint
echo "🏥 Testing health endpoint..."
HEALTH_RESPONSE=$(docker exec fluently-llm-test curl -s -f http://localhost:8003/health)

if [ $? -eq 0 ]; then
    echo "✅ Health endpoint accessible"
    echo "Response: $HEALTH_RESPONSE"
else
    echo "❌ Health endpoint failed"
    docker logs fluently-llm-test
    docker stop fluently-llm-test >/dev/null 2>&1
    exit 1
fi

# Stop container
echo "🛑 Stopping test container..."
docker stop fluently-llm-test >/dev/null 2>&1

# Test with real API keys if provided
if [ ! -z "$GROQ_API_KEYS" ] || [ ! -z "$GEMINI_API_KEYS" ]; then
    echo "🔑 Testing with real API keys..."
    
    docker run --rm -d \
        -e GROQ_API_KEYS="$GROQ_API_KEYS" \
        -e GEMINI_API_KEYS="$GEMINI_API_KEYS" \
        --name fluently-llm-test-real \
        fluently-llm-api:test
    
    # Wait for initialization
    sleep 15
    
    # Test chat endpoint
    echo "💬 Testing chat endpoint..."
    CHAT_RESPONSE=$(docker exec fluently-llm-test-real curl -s -f \
        -X POST "http://localhost:8003/chat" \
        -H "Content-Type: application/json" \
        -d '{"messages": [{"role": "user", "content": "Say hello"}], "model_type": "fast"}')
    
    if [ $? -eq 0 ]; then
        echo "✅ Chat endpoint works with real API keys"
        echo "Response: $CHAT_RESPONSE"
    else
        echo "❌ Chat endpoint failed with real API keys"
        docker logs fluently-llm-test-real
    fi
    
    docker stop fluently-llm-test-real >/dev/null 2>&1
else
    echo "⚠️  No real API keys provided, skipping chat test"
fi

echo "🎉 Test completed successfully!"
