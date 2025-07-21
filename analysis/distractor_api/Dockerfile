# Dockerfile for FastAPI Distractor API (CPU-Only Optimized)
FROM python:3.11-slim AS base

# Set environment variables
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONPATH=/app \
    PIP_DISABLE_PIP_VERSION_CHECK=1

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user for security
RUN useradd --create-home --shell /bin/bash app

WORKDIR /app

# ============== Dependency Layer (Cached) ==============
FROM base as dependencies

# Copy only requirements first (for better layer caching)
COPY distractor_api/requirements.txt ./

# Install dependencies with caching enabled
# This layer will be cached as long as requirements.txt doesn't change
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install --find-links https://download.pytorch.org/whl/cpu -r requirements.txt

# ============== Application Layer ==============
FROM dependencies as final

# Copy pyproject.toml (separate layer for metadata)
COPY distractor_api/pyproject.toml ./

# Copy application code and bert module (separate layers for better caching)
COPY --chown=app:app bert/ ./bert
COPY --chown=app:app distractor_api/ .

# Create cache directory for models and set permissions
RUN mkdir -p /app/.cache/huggingface && \
    chown -R app:app /app/.cache

# Switch to non-root user
USER app

# Download models at build time to cache them in the image
# This prevents download time during first API call
RUN python -c "from transformers import BertTokenizerFast, BertForMaskedLM; BertTokenizerFast.from_pretrained('bert-base-uncased'); BertForMaskedLM.from_pretrained('bert-base-uncased')" || echo "Model download optional - will download at runtime if needed"

# Download spaCy model at build time
RUN python -m spacy download en_core_web_sm

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5m --retries=3 \
    CMD curl -f http://localhost:8001/health || exit 1

# Expose port
EXPOSE 8001

# Run the application
CMD ["python", "-m", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8001", "--workers", "1"]