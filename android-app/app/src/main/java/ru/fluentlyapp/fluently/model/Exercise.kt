package ru.fluentlyapp.fluently.model

sealed interface Exercise : LessonComponent {
    val isAnswered: Boolean

    data class NewWord(
        val word: String,
        val translation: String,
        val phoneticTranscription: String, // like /ˈkɑːn.ʃəs.nəs/
        val doesUserKnow: Boolean?,
        val examples: List<Pair<String, String>> // Pair(sentence, translation)
    ) : Exercise {
        override val isAnswered = doesUserKnow == null
    }

    data class ChooseTranslation(
        val word: String,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
    }

    data class FillTheGap(
        val sentence: List<String>,
        val answerVariants: List<String>,
        val correctVariant: Int,
        val selectedVariant: Int?
    ) : Exercise {
        override val isAnswered = selectedVariant != null
    }

    data class InputWord(
        val translation: String,
        val correctAnswer: String,
        val inputtedWord: String?
    ) : Exercise {
        override val isAnswered = inputtedWord != null
    }
}
