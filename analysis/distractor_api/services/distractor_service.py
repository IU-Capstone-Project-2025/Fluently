import asyncio
import random
from typing import List, Optional
import logging
import sys
from pathlib import Path

# Add parent directory to path to import bert module
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from bert.distractor_generator import DistractorGenerator

logger = logging.getLogger(__name__)


class DistractorService:
    """High-performance distractor generation service using existing DistractorGenerator"""

    def __init__(self):
        self.generator: Optional[DistractorGenerator] = None
        self.initialized = False

    async def initialize(self):
        """Initialize the model asynchronously"""
        if self.initialized:
            return

        try:
            # Load the DistractorGenerator in a separate thread to avoid blocking
            loop = asyncio.get_event_loop()
            await loop.run_in_executor(None, self._load_model)
            self.initialized = True
            logger.info("BERT DistractorGenerator loaded successfully")
        except Exception as e:
            logger.error(f"Failed to load BERT model: {e}")
            raise RuntimeError(f"Model initialization failed: {e}")

    def _load_model(self):
        """Load the DistractorGenerator"""
        self.generator = DistractorGenerator()

        # Warm up the model with a dummy input for better first-request latency
        try:
            self.generator.generate_distractors("The cat is sleeping", "cat", 1)
            logger.info("Model warmed up successfully")
        except Exception as e:
            logger.warning(f"Model warmup failed: {e}")

    async def generate_distractors(
        self, sentence: str, target_word: str, num_distractors: int = 3
    ) -> List[str]:
        """
        Generate distractors for a target word in a sentence using existing DistractorGenerator

        Args:
            sentence: The input sentence
            target_word: The word to generate distractors for
            num_distractors: Number of distractors to generate (default: 3)

        Returns:
            List of distractor words including the original word
        """
        if not self.initialized or not self.generator:
            raise RuntimeError("Service not initialized")

        try:
            # Use the existing DistractorGenerator in a separate thread to avoid blocking
            loop = asyncio.get_event_loop()
            distractors = await loop.run_in_executor(
                None,
                self.generator.generate_distractors,
                sentence,
                target_word,
                num_distractors  # Now properly passing the num_distractors parameter
            )

            # Handle empty results from generator
            if not distractors:
                logger.warning(f"No distractors generated for word '{target_word}' in sentence: {sentence}")
                # Return just the original word if no distractors were generated
                return [target_word]

            # Always include the original word and shuffle for randomness
            all_options = [target_word] + distractors
            random.shuffle(all_options)

            return all_options

        except Exception as e:
            logger.error(f"Error generating distractors for word '{target_word}': {e}")
            # Fallback to original word only
            return [target_word]
