import os
import asyncio
import logging
from groq import Groq, APIError
from google import generativeai as genai
from google.api_core import exceptions as google_exceptions
from datetime import datetime, timedelta

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("AIService")

class AIService:
    def __init__(self):
        self.providers = self._init_providers()
        self.fallback_strategy = {
            "groq": "gemini",
            "gemini": "groq"
        }
        self.key_rotation = {provider: 0 for provider in self.providers}
        self.failures = {provider: {} for provider in self.providers}  # {provider: {key: failure_count}}
        self.cooldown_period = timedelta(minutes=5)
        self.request_timeout = 15  # секунд

    def _init_providers(self):
        """Инициализация провайдеров и их API-ключей"""
        # Helper function to safely split API keys
        def get_api_keys(env_var):
            keys_str = os.getenv(env_var, "").strip()
            if not keys_str:
                return []
            return [key.strip() for key in keys_str.split(",") if key.strip()]
        
        return {
            "groq": {
                "keys": get_api_keys("GROQ_API_KEYS"),
                "models": {
                    "fast": "llama-3-8b-8192",
                    "balanced": "llama-3.3-70b-versatile"
                },
                "client": None
            },
            "gemini": {
                "keys": get_api_keys("GEMINI_API_KEYS"),
                "models": {
                    "fast": "gemini-2.0-flash",
                    "balanced": "gemini-2.5-flash"
                },
                "client": None
            }
        }

    async def initialize(self):
        """Асинхронная инициализация клиентов"""
        await asyncio.to_thread(self._sync_initialize)

    def _sync_initialize(self):
        """Синхронная инициализация клиентов"""
        # Initialize Groq client with first available key
        groq_keys = self.providers["groq"]["keys"]
        if groq_keys:
            try:
                self.providers["groq"]["client"] = Groq(api_key=groq_keys[0])
                logger.info(f"Groq client initialized with {len(groq_keys)} key(s)")
            except Exception as e:
                logger.error(f"Failed to initialize Groq client: {e}")
                logger.warning("Groq client initialization failed, continuing with Gemini only")
        else:
            logger.warning("No Groq API keys found in environment")
            
        # Initialize Gemini with first available key
        gemini_keys = self.providers["gemini"]["keys"]
        if gemini_keys:
            try:
                genai.configure(api_key=gemini_keys[0])
                logger.info(f"Gemini client initialized with {len(gemini_keys)} key(s)")
            except Exception as e:
                logger.error(f"Failed to initialize Gemini client: {e}")
                logger.warning("Gemini client initialization failed")
        else:
            logger.warning("No Gemini API keys found in environment")
            
        # Check if at least one provider is available
        if not any(self.providers[provider]["keys"] for provider in self.providers):
            raise RuntimeError("No valid API keys found for any provider (Groq, Gemini)")
            
        logger.info("AI service initialization completed")

    def _get_next_key(self, provider):
        """Получение следующего API-ключа с ротацией"""
        keys = self.providers[provider]["keys"]
        if not keys:
            return None
        
        idx = self.key_rotation[provider] % len(keys)
        self.key_rotation[provider] = (idx + 1) % len(keys)
        return keys[idx]

    def _mark_failure(self, provider, key):
        """Отметить неудачный запрос для ключа"""
        if key not in self.failures[provider]:
            self.failures[provider][key] = {"count": 0, "last_failure": datetime.min}
        
        self.failures[provider][key]["count"] += 1
        self.failures[provider][key]["last_failure"] = datetime.now()
        logger.warning(f"Failure #{self.failures[provider][key]['count']} for {provider} key: {key[-6:]}")

    def _is_key_cooldown(self, provider, key):
        """Проверить, находится ли ключ в режиме охлаждения"""
        if key not in self.failures[provider]:
            return False
        
        failure_info = self.failures[provider][key]
        time_since_failure = datetime.now() - failure_info["last_failure"]
        
        # Экспоненциальное увеличение времени охлаждения
        cooldown = min(self.cooldown_period * (2 ** (failure_info["count"] - 1)), timedelta(hours=1))
        return time_since_failure < cooldown

    async def chat_completion(self, messages, model_type="balanced", **kwargs):
        """
        Асинхронное получение ответа от ИИ-модели
        :param messages: История сообщений в формате [{"role": "user", "content": "..."}]
        :param model_type: Тип модели ("fast", "balanced")
        :return: Ответ ИИ
        """
        providers_order = ["groq", "gemini"]
        response = None
        
        for provider in providers_order:
            key = self._get_next_key(provider)
            if not key:
                continue
                
            if self._is_key_cooldown(provider, key):
                logger.info(f"Skipping {provider} key in cooldown: {key[-6:]}")
                continue
                
            try:
                if provider == "groq":
                    response = await self._groq_request(key, messages, model_type, **kwargs)
                elif provider == "gemini":
                    response = await self._gemini_request(key, messages, model_type, **kwargs)
                
                if response:
                    return response
                    
            except Exception as e:
                self._handle_error(provider, key, e)
        
        raise RuntimeError("All providers failed")

    async def _groq_request(self, key, messages, model_type, **kwargs):
        """Запрос к Groq API"""
        model_name = self.providers["groq"]["models"][model_type]
        
        try:
            client = Groq(api_key=key)
        except Exception as e:
            logger.error(f"Failed to create Groq client: {e}")
            return None
        
        try:
            response = await asyncio.wait_for(
                asyncio.to_thread(
                    client.chat.completions.create,
                    messages=messages,
                    model=model_name,
                    **kwargs
                ),
                timeout=self.request_timeout
            )
            return response.choices[0].message.content
        except APIError as e:
            if e.status_code == 429:
                logger.warning(f"Groq rate limit exceeded for key: {key[-6:]}")
                return None
            raise

    async def _gemini_request(self, key, messages, model_type, **kwargs):
        """Запрос к Gemini API"""
        model_name = self.providers["gemini"]["models"][model_type]
        genai.configure(api_key=key)
        model = genai.GenerativeModel(model_name)
        
        # Преобразование формата сообщений
        gemini_messages = self._convert_to_gemini_format(messages)
        
        try:
            response = await asyncio.wait_for(
                model.generate_content_async(gemini_messages, **kwargs),
                timeout=self.request_timeout
            )
            return response.text
        except google_exceptions.ResourceExhausted:
            logger.warning(f"Gemini quota exceeded for key: {key[-6:]}")
            return None

    def _convert_to_gemini_format(self, messages):
        """Конвертация формата сообщений в совместимый с Gemini"""
        gemini_messages = []
        
        for msg in messages:
            role = "user" if msg["role"] in ["user", "system"] else "model"
            gemini_messages.append({"role": role, "parts": [{"text": msg["content"]}]})
        
        return gemini_messages

    def _handle_error(self, provider, key, error):
        """Обработка ошибок и обновление состояния ключей"""
        error_type = type(error).__name__
        logger.error(f"{provider} error ({error_type}) with key {key[-6:]}: {str(error)}")
        
        # Отметить неудачу
        self._mark_failure(provider, key)
        
        # Если это критическая ошибка аутентификации
        if "401" in str(error) or "403" in str(error):
            logger.error(f"Invalid API key detected for {provider}: {key[-6:]}")
            if key in self.providers[provider]["keys"]:
                self.providers[provider]["keys"].remove(key)