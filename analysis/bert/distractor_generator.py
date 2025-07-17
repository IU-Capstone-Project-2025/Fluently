from transformers import BertTokenizerFast, BertForMaskedLM
import torch
import re
import random
import spacy
from collections import OrderedDict

class DistractorGenerator:
    def __init__(self):
        # Загрузка модели BERT
        self.model_name = "bert-base-uncased"
        self.tokenizer = BertTokenizerFast.from_pretrained(self.model_name)
        self.model = BertForMaskedLM.from_pretrained(self.model_name)
        self.model.eval()

        # Загрузка модели spaCy для английского языка
        try:
            self.nlp = spacy.load("en_core_web_sm")
        except OSError:
            # Автоматическая установка модели при первом запуске
            import subprocess
            import sys
            subprocess.run([sys.executable, "-m", "spacy", "download", "en_core_web_sm"], check=True)
            self.nlp = spacy.load("en_core_web_sm")

    def is_valid_distractor(self, token: str, target_lemma: str) -> bool:
        """Проверяет, является ли токен валидным дистрактором"""
        # Проверка длины
        if len(token) < 2:
            return False

        # Проверка на содержание только букв
        if not token.isalpha():
            return False

        # Проверка на части слов (BERT-specific)
        if token.startswith("##"):
            return False

        # Проверка на совпадение с целевой леммой
        try:
            token_lemma = self.nlp(token)[0].lemma_.lower()
        except:
            token_lemma = token.lower()

        return token_lemma != target_lemma

    def generate_distractors(self, sentence: str, target_word: str, num_distractors: int = 3) -> list:
        # Обработка предложения с помощью spaCy
        doc = self.nlp(sentence)
        target_lemma = self.nlp(target_word)[0].lemma_.lower()

        # Поиск слова для маскирования по лемме
        word_to_mask = None
        for token in doc:
            if token.lemma_.lower() == target_lemma:
                word_to_mask = token.text
                break

        if word_to_mask is None:
            return []

        # Создаем маскированное предложение
        masked_sentence = re.sub(
            rf"\b{re.escape(word_to_mask)}\b",
            self.tokenizer.mask_token,
            sentence,
            count=1,
            flags=re.IGNORECASE
        )

        if masked_sentence == sentence:
            return []

        # Получаем предсказания от BERT
        inputs = self.tokenizer(masked_sentence, return_tensors="pt")
        mask_token_id = self.tokenizer.mask_token_id
        mask_positions = torch.where(inputs["input_ids"][0] == mask_token_id)[0]

        distractors = []
        with torch.no_grad():
            outputs = self.model(**inputs)

        for pos in mask_positions:
            logits = outputs.logits[0, pos]
            probabilities = torch.softmax(logits, dim=-1)

            # Берем топ-50 кандидатов
            top_k = 50
            top_probs, top_indices = torch.topk(probabilities, top_k)

            candidates = []
            for idx in top_indices:
                token = self.tokenizer.decode(idx).strip()

                # Пропускаем невалидные варианты
                if not self.is_valid_distractor(token, target_lemma):
                    continue

                candidates.append(token)

            # Удаляем дубликаты с сохранением порядка
            candidates = list(OrderedDict.fromkeys(candidates))

            # Выбираем "средние" по вероятности варианты
            if candidates:
                start_idx = min(5, len(candidates))
                end_idx = min(start_idx + 10, len(candidates))

                if end_idx > start_idx:
                    mid_candidates = candidates[start_idx:end_idx]
                    distractors = random.sample(mid_candidates, min(num_distractors, len(mid_candidates)))
                else:
                    distractors = random.sample(candidates, min(num_distractors, len(candidates)))
                break

        return distractors[:num_distractors]

# Пример использования
if __name__ == "__main__":
    generator = DistractorGenerator()
    examples = [
        ("The cat caught the mouse in the kitchen", "mouse"),
        ("She opened the window to breathe fresh air", "window"),
        ("The programmer fixed all bugs in the code", "bugs"),
        ("They are playing football in the yard", "playing"),
        ("Birds fly south for the winter", "bird"),
        ("He played football yesterday", "play"),
        ("These mice are tiny", "mouse"),
        ("Children are playing outside", "child"),
        ("Women in business", "woman"),
        ("The geese are flying south", "goose"),
        ("His feet are cold", "foot"),
        ("She has a new pair of teeth", "tooth"),
        ("A cat is on the mat", "cat")  # Проверка на однобуквенные слова
    ]

    for sentence, word in examples:
        distractors = generator.generate_distractors(sentence, word)
        print(f"Предложение: {sentence}")
        print(f"Целевое слово: {word}")
        print(f"Дистракторы: {', '.join(distractors) if distractors else 'No distractors found'}\n")
