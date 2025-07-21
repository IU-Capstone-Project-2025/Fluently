package ru.fluentlyapp.fluently.common.model

import androidx.browser.customtabs.CustomTabsService
import kotlinx.serialization.Serializable

sealed interface LessonComponent {
    var id: Int
}

sealed interface Exercise : LessonComponent {
    val isAnswered: Boolean
    val isCorrect: Boolean?

    @Serializable
    data class NewWord(
        override var id: Int = -1,
        val wordId: String = "",
        val word: String,
        val translation: String,
        val phoneticTranscription: String,
        val doesUserKnow: Boolean?,
        val examples: List<Pair<String, String>>
    ) : Exercise {
        override val isAnswered = doesUserKnow != null
        override val isCorrect: Boolean?
            get() = true
    }

    @Serializable
    data class ChooseTranslation(
        override var id: Int = -1,
        val wordId: String = "",
        val word: String,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
        override val isCorrect: Boolean?
            get() = if (isAnswered) selectedVariant == correctVariant else null
    }

    @Serializable
    data class FillTheGap(
        override var id: Int = -1,
        val wordId: String = "",
        val sentence: List<String>,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
        override val isCorrect: Boolean?
            get() = if (isAnswered) selectedVariant == correctVariant else null
    }

    @Serializable
    data class InputWord(
        val translation: String,
        val wordId: String = "",
        val correctAnswer: String,
        val inputtedWord: String?,
        override var id: Int = -1,
    ) : Exercise {
        override val isAnswered = inputtedWord != null
        override val isCorrect: Boolean?
            get() = if (isAnswered) inputtedWord?.trim() == correctAnswer.trim() else null
    }
}

@Serializable
data class Dialog(
    val messages: List<Message>,
    val isFinished: Boolean,
    override var id: Int = -1
) : LessonComponent {
    @Serializable
    data class Message(
        val messageId: Long,
        val text: String,
        val fromUser: Boolean,
    )
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
    class Finish(override var id: Int = -1) : Decoration

    @Serializable
    class LearningPartComplete(override var id: Int = -1) : Decoration

}
