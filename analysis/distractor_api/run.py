#!/usr/bin/env python3
"""
Standalone runner for the Distractor API service.
This script handles the Python path setup and starts the FastAPI server.
"""

import sys
from pathlib import Path

# Add the parent directory to Python path for imports (for standalone usage)
current_dir = Path(__file__).parent.absolute()
parent_dir = current_dir.parent
if str(parent_dir) not in sys.path:
    sys.path.insert(0, str(parent_dir))

from distractor_api.main import app


def main():
    """Main entry point for Poetry script"""
    import uvicorn

    print("Starting Distractor API service...")
    print(f"Current directory: {current_dir}")
    print(f"Python path: {sys.path[:3]}...")

    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8001,
        reload=False,  # Disable reload for better performance
        workers=1,  # Single worker to avoid model loading overhead
        log_level="info",
    )


if __name__ == "__main__":
    main()
