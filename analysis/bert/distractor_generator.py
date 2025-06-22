from transformers import BertTokenizerFast, BertForMaskedLM
import torch
import re
import random

class DistractorGenerator:
    def __init__(self):
        # Загрузка модели BERT
        self.model_name = "bert-base-uncased"
        self.tokenizer = BertTokenizerFast.from_pretrained(self.model_name)
        self.model = BertForMaskedLM.from_pretrained(self.model_name)
        self.model.eval()

    def generate_distractors(self, sentence: str, target_word: str, num_distractors: int = 3) -> list:
        """
        Генерирует дистракторы для заданного слова в предложении

        Параметры:
        sentence: Исходное предложение
        target_word: Слово для замены
        num_distractors: Количество дистракторов для генерации

        Возвращает:
        list: Список дистракторов
        """
        # Создаем маскированное предложение
        masked_sentence = re.sub(
            rf"\b{re.escape(target_word)}\b",
            self.tokenizer.mask_token,
            sentence,
            count=1,
            flags=re.IGNORECASE
        )

        # Проверяем замену
        if masked_sentence == sentence:
            return []

        # Токенизация
        inputs = self.tokenizer(masked_sentence, return_tensors="pt")
        mask_token_id = self.tokenizer.mask_token_id

        # Находим позицию маски
        mask_positions = torch.where(inputs["input_ids"][0] == mask_token_id)[0]

        # Получаем предсказания
        with torch.no_grad():
            outputs = self.model(**inputs)

        # Собираем результаты
        distractors = []
        target_word_lower = target_word.lower()

        for pos in mask_positions:
            logits = outputs.logits[0, pos]
            probabilities = torch.softmax(logits, dim=-1)

            # Берем топ-50 кандидатов
            top_k = 50
            top_probs, top_indices = torch.topk(probabilities, top_k)

            # Фильтруем и собираем кандидатов
            candidates = []
            for i, idx in enumerate(top_indices):
                token = self.tokenizer.decode(idx).strip()
                token_lower = token.lower()

                # Пропускаем невалидные варианты
                if (not token or
                    token.startswith("##") or
                    not token.isalpha() or
                    token_lower == target_word_lower):
                    continue

                candidates.append(token)

            # Удаляем дубликаты
            candidates = list(dict.fromkeys(candidates))

            # Выбираем "слегка неподходящие" варианты
            if len(candidates) > 10:
                # Пропускаем первые 5 самых очевидных
                mid_candidates = candidates[5:15]

                # Выбираем случайные из "серединки"
                distractors = random.sample(mid_candidates, min(num_distractors, len(mid_candidates)))
                break
            elif candidates:
                distractors = random.sample(candidates, min(num_distractors, len(candidates)))

        return distractors[:num_distractors]

# Пример использования
if __name__ == "__main__":
    generator = DistractorGenerator()

    examples = [
        ("The cat caught the mouse in the kitchen", "mouse"),
        ("She opened the window to breathe fresh air", "window"),
        ("The programmer fixed all bugs in the code", "bugs"),
        ("They are playing football in the yard", "playing"),
        ("Birds fly south for the winter", "fly")
    ]

    for sentence, word in examples:
        distractors = generator.generate_distractors(sentence, word)
        print(f"Предложение: {sentence}")
        print(f"Целевое слово: {word}")
        print(f"Дистракторы: {', '.join(distractors)}\n")