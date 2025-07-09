package ru.fluentlyapp.fluently.common.model

import kotlinx.serialization.Serializable

sealed interface LessonComponent {
    var id: Int
}

sealed interface Exercise : LessonComponent {
    val isAnswered: Boolean

    @Serializable
    data class NewWord(
        override var id: Int = -1,
        val word: String,
        val translation: String,
        val phoneticTranscription: String,
        val doesUserKnow: Boolean?,
        val examples: List<Pair<String, String>>
    ) : Exercise {
        override val isAnswered = doesUserKnow != null
    }

    @Serializable
    data class ChooseTranslation(
        override var id: Int = -1,
        val word: String,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
    }

    @Serializable
    data class FillTheGap(
        override var id: Int = -1,
        val sentence: List<String>,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
    }

    @Serializable
    data class InputWord(
        val translation: String,
        val correctAnswer: String,
        val inputtedWord: String?,
        override var id: Int = -1,
    ) : Exercise {
        override val isAnswered = inputtedWord != null
    }
}

@Serializable
sealed interface Decoration : LessonComponent {
    @Serializable
    class Loading(override var id: Int = -1) : Decoration

    @Serializable
    data class Onboarding(
        val wordsToBeLearned: Int,
        val featuredExercises: Int,
        override var id: Int = -1,
    ) : Decoration

    @Serializable
    class Finish(override var id: Int = -1): Decoration
}
