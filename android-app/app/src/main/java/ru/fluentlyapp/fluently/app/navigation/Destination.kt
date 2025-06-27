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
    data class LessonScreen(val lessonId: String)
}