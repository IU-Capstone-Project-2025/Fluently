# Dockerfile for FastAPI Distractor API (CPU-Only Optimized)
FROM python:3.11-slim

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONPATH=/app \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user for security
RUN useradd --create-home --shell /bin/bash app

WORKDIR /app

# Copy requirements and install dependencies
COPY distractor_api/pyproject.toml ./

# Install Poetry and export requirements
RUN pip install poetry==1.7.1 && \
    poetry export -f requirements.txt --output requirements.txt --without-hashes

# Install CPU-only PyTorch first to avoid CUDA dependencies
RUN pip install torch==2.1.0+cpu torchvision==0.16.0+cpu torchaudio==0.16.0+cpu \
    --index-url https://download.pytorch.org/whl/cpu

# Install remaining dependencies (excluding torch since we already installed CPU version)
RUN pip install -r requirements.txt --no-deps && \
    pip install transformers tokenizers numpy pandas nltk scikit-learn \
    fastapi uvicorn[standard] pydantic python-multipart requests && \
    pip uninstall poetry -y

# Copy application code and bert module
COPY --chown=app:app distractor_api/ .
COPY --chown=app:app bert/ ./bert

# Create cache directory for models and set permissions
RUN mkdir -p /app/.cache/huggingface && \
    chown -R app:app /app/.cache

# Switch to non-root user
USER app

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5m --retries=3 \
    CMD curl -f http://localhost:8001/health || exit 1

# Expose port
EXPOSE 8001

# Run the application
CMD ["python", "-m", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8001", "--workers", "1"] 