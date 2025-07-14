package ru.fluentlyapp.fluently.app.navigation

import kotlinx.serialization.Serializable

object Destination {
    @Serializable
    object LaunchScreen

    @Serializable
    object LoginScreen

    @Serializable
    object HomeScreen

    @Serializable
    object LessonScreen

    @Serializable
    class WordsProgress(val progressType: WordsProgressType) {
        enum class WordsProgressType {
            LEARNED, IN_PROGRESS
        }
    }
}
