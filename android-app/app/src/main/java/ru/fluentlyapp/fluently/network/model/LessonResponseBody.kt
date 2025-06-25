package ru.fluentlyapp.fluently.network.model

import kotlinx.serialization.Serializable

@Serializable
data class LessonResponseBody(
    val lesson: LessonApiModel,
    val cards: List<CardApiModel>,
    val sync: Sync? = null
)

@Serializable
data class LessonApiModel(
    val lesson_id: String,
    val user_id: String,
    val started_at: String,
    val words_per_lesson: Int,
    val total_words: Int
)

@Serializable
data class CardApiModel(
    val word_id: String,
    val word: String,
    val translation: String,
    val transcription: String,
    val cefr_level: String,
    val is_new: Boolean,
    val topic: String,
    val subtopic: String,
    val sentences: List<SentenceApiModel>,
    val exercise: ExerciseApiModel
)

@Serializable
data class SentenceApiModel(
    val sentence_id: String,
    val text: String,
    val translation: String
)

@Serializable
data class ExerciseApiModel(
    val exercise_id: String,
    val type: String,
    val data: ExerciseDataApiModel
) {
    object ExerciseType {
        const val TRANSLATE_RU_TO_EN = "translate_ru_to_en"
        const val TRANSLATE_EN_TO_RU = "translate_en_to_ru"
        const val PICK_OPTIONS_SENTENCE = "pick_option_sentence"
    }
}

@Serializable
data class ExerciseDataApiModel(
    // All optional fields to handle different exercise types
    val text: String? = null,
    val correct_answer: String? = null,
    val pick_options: List<String>? = null,
    val translation: String? = null,
    val sentence_id: String? = null,
    val template: String? = null
)

@Serializable
data class Sync(
    val dirty: Boolean,
    val last_synced_at: String
)
