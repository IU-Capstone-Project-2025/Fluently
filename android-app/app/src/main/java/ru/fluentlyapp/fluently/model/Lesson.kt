package ru.fluentlyapp.fluently.model

data class Lesson(
    val lessonId: String,
    val components: List<LessonComponent>,
    val currentLessonComponentIndex: Int = 0
)