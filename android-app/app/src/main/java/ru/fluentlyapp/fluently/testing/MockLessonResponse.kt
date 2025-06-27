package ru.fluentlyapp.fluently.testing

import ru.fluentlyapp.fluently.network.model.CardApiModel
import ru.fluentlyapp.fluently.network.model.ExerciseApiModel
import ru.fluentlyapp.fluently.network.model.ExerciseDataApiModel
import ru.fluentlyapp.fluently.network.model.LessonApiModel
import ru.fluentlyapp.fluently.network.model.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.SentenceApiModel
import ru.fluentlyapp.fluently.network.model.Sync

val mockLessonResponse = LessonResponseBody(
    lesson = LessonApiModel(
        lesson_id = "3f7c9f5e-8c1e-4b2e-9f57-23eac9d8b123",
        user_id = "9d8e7f6c-5b4a-3f2e-1d0c-9b8a7f6e5d4c",
        started_at = "2025-06-24T10:30:00Z",
        words_per_lesson = 10,
        total_words = 30
    ),
    cards = listOf(
        CardApiModel(
            word_id = "dcba4321-8765-09ba-fedc-222222222222",
            word = "car",
            translation = "машина",
            transcription = "[kɑːr]",
            cefr_level = "A1",
            is_new = false,
            topic = "Transport",
            subtopic = "Vehicles",
            sentences = listOf(
                SentenceApiModel(
                    sentence_id = "33334444-5555-6666-7777-888899990000",
                    text = "He drives a car to work every day.",
                    translation = "Он ездит на машине на работу каждый день."
                )
            ),
            exercise = ExerciseApiModel(
                exercise_id = "00009999-8888-7777-6666-555544443333",
                type = "translate_ru_to_en",
                data = ExerciseDataApiModel(
                    text = "машина",
                    correct_answer = "car",
                    pick_options = listOf("bus", "bike", "train")
                )
            )
        ),
        CardApiModel(
            word_id = "efgh5678-1234-90ab-cdef-333333333333",
            word = "apple",
            translation = "яблоко",
            transcription = "[ˈæpəl]",
            cefr_level = "A1",
            is_new = true,
            topic = "Food",
            subtopic = "Fruits",
            sentences = listOf(
                SentenceApiModel(
                    sentence_id = "44445555-6666-7777-8888-999900001111",
                    text = "I ate an apple for breakfast.",
                    translation = "Я съел яблоко на завтрак."
                )
            ),
            exercise = ExerciseApiModel(
                exercise_id = "aaaabbbb-cccc-dddd-eeee-ffff00001111",
                type = "write_word_from_translation",
                data = ExerciseDataApiModel(
                    translation = "яблоко",
                    correct_answer = "apple"
                )
            )
        ),
        CardApiModel(
            word_id = "ijkl9012-3456-78ab-cdef-444444444444",
            word = "house",
            translation = "дом",
            transcription = "[haʊs]",
            cefr_level = "A1",
            is_new = true,
            topic = "Home",
            subtopic = "Buildings",
            sentences = listOf(
                SentenceApiModel(
                    sentence_id = "55556666-7777-8888-9999-000011112222",
                    text = "I built a house last year.",
                    translation = "Я построил дом в прошлом году."
                )
            ),
            exercise = ExerciseApiModel(
                exercise_id = "ccccdddd-eeee-ffff-1111-222233334444",
                type = "pick_option_sentence",
                data = ExerciseDataApiModel(
                    sentence_id = "55556666-7777-8888-9999-000011112222",
                    template = "I built a _ last year.",
                    correct_answer = "house",
                    pick_options = listOf("tree", "car", "computer")
                )
            )
        )
    ),
    sync = Sync(
        dirty = true,
        last_synced_at = "2025-06-24T10:35:45Z"
    )
)
