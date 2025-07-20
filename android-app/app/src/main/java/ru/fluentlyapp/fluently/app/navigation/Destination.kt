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
    class WordsProgress(val isLearning: Boolean)

    @Serializable
    object CalendarScreen

    @Serializable
    object OnboardingScreen
}
