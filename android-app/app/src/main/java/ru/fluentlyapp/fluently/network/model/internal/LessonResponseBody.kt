package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class LessonResponseBody(
    val lesson: LessonApiModel,
    val cards: List<CardApiModel>
)

@Serializable
data class LessonApiModel(
    val started_at: String,
    val words_per_lesson: Int,
    val total_words: Int,
    val cefr_level: String
)

@Serializable
data class CardApiModel(
    val word_id: String,
    val word: String,
    val translation: String,
    val topic: String,
    val subtopic: String,
    val sentences: List<SentenceApiModel>,
    val exercise: ExerciseApiModel
)

@Serializable
data class SentenceApiModel(
    val text: String,
    val translation: String
)

@Serializable
data class ExerciseApiModel(
    val type: String,
    val data: ExerciseDataApiModel
) {
    object ExerciseType {
        const val TRANSLATE_RU_TO_EN = "translate_ru_to_en"
        const val TRANSLATE_EN_TO_RU = "translate_en_to_ru"
        const val PICK_OPTION_SENTENCE = "pick_option_sentence"
        const val WRITE_WORD_FROM_TRANSLATION = "write_word_from_translation"
    }
}

@Serializable
data class ExerciseDataApiModel(
    val text: String? = null,
    val correct_answer: String? = null,
    val pick_options: List<String>? = null,
    val translation: String? = null,
    val template: String? = null
)