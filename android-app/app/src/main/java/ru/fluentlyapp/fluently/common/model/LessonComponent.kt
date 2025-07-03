package ru.fluentlyapp.fluently.common.model

import kotlinx.serialization.Serializable

@Serializable
sealed class LessonComponent {
    var id: Int = -1
}

@Serializable
sealed class Exercise : LessonComponent() {
    abstract val isAnswered: Boolean

    @Serializable
    data class NewWord(
        val word: String,
        val translation: String,
        val phoneticTranscription: String, // like /ˈkɑːn.ʃəs.nəs/
        val doesUserKnow: Boolean?,
        val examples: List<Pair<String, String>> // Pair(sentence, translation)
    ) : Exercise() {
        override val isAnswered = doesUserKnow != null
    }

    @Serializable
    data class ChooseTranslation(
        val word: String,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise() {
        override val isAnswered = selectedVariant != null
    }

    @Serializable
    data class FillTheGap(
        val sentence: List<String>,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise() {
        override val isAnswered = selectedVariant != null
    }

    @Serializable
    data class InputWord(
        val translation: String,
        val correctAnswer: String,
        val inputtedWord: String?
    ) : Exercise() {
        override val isAnswered = inputtedWord != null
    }
}

@Serializable
sealed class Decoration : LessonComponent() {
    @Serializable
    object Loading : Decoration()
}