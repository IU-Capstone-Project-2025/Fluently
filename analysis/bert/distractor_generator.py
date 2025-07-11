from transformers import BertTokenizerFast, BertForMaskedLM
import torch
import re
import random

class DistractorGenerator:
    def __init__(self):
        self.model_name = "bert-base-uncased"
        self.tokenizer = BertTokenizerFast.from_pretrained(self.model_name)
        self.model = BertForMaskedLM.from_pretrained(self.model_name)
        self.model.eval()

        # Простой словарь для ручной лемматизации
        self.lemmatization_map = {
            "mice": "mouse",
            "geese": "goose",
            "feet": "foot",
            "teeth": "tooth",
            "knives": "knife",
            "playing": "play",
            "played": "play",
            "plays": "play",
            "broke": "break",
            "broken": "break",
            "began": "begin",
            "begun": "begin",
            "birds": "bird",
            "cats": "cat",
            "dogs": "dog",
            "children": "child",
            "men": "man",
            "women": "woman"
        }

    def custom_lemmatize(self, word: str) -> str:
        """Простая ручная лемматизация для английского языка"""
        word_lower = word.lower()

        # Сначала проверяем исключения
        if word_lower in self.lemmatization_map:
            return self.lemmatization_map[word_lower]

        # Правила для регулярных форм
        if word_lower.endswith("ies") and len(word_lower) > 3:
            return word_lower[:-3] + "y"
        if word_lower.endswith("es") and len(word_lower) > 2:
            return word_lower[:-2]
        if word_lower.endswith("s") and len(word_lower) > 1:
            return word_lower[:-1]
        if word_lower.endswith("ing") and len(word_lower) > 4:
            return word_lower[:-3]
        if word_lower.endswith("ed") and len(word_lower) > 3:
            return word_lower[:-2]

        return word_lower

    def generate_distractors(self, sentence: str, target_word: str, num_distractors: int = 3) -> list:
        # Лемматизируем целевое слово
        target_lemma = self.custom_lemmatize(target_word).lower()

        # Ищем слово для маскирования
        words = re.findall(r'\b\w+\b', sentence)
        word_to_mask = None

        for word in words:
            word_lemma = self.custom_lemmatize(word).lower()
            if word_lemma == target_lemma:
                word_to_mask = word
                break

        if word_to_mask is None:
            return []

        # Создаем маскированное предложение
        masked_sentence = re.sub(
            rf"\b{re.escape(word_to_mask)}\b",
            self.tokenizer.mask_token,
            sentence,
            count=1
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
                token_lower = token.lower()

                # Пропускаем невалидные варианты
                if (not token
                    or token.startswith("##")
                    or not token.isalpha()
                    or self.custom_lemmatize(token) == target_lemma):
                    continue

                candidates.append(token)

            # Удаляем дубликаты
            candidates = list(dict.fromkeys(candidates))

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
        ("Women in business", "woman")
    ]

    for sentence, word in examples:
        distractors = generator.generate_distractors(sentence, word)
        print(f"Предложение: {sentence}")
        print(f"Целевое слово: {word}")
        print(f"Дистракторы: {', '.join(distractors)}\n")