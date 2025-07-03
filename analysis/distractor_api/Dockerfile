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

# Copy requirements and pyproject files
COPY distractor_api/requirements.txt ./
COPY distractor_api/pyproject.toml ./

# Install dependencies directly from requirements.txt (much faster than Poetry export)
RUN pip install --no-cache-dir -r requirements.txt

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