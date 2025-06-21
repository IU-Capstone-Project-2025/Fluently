package ru.fluentlyapp.fluently.model

data class Lesson(
    val lessonId: String,
    val exercises: List<Exercise>
)