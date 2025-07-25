package ru.fluentlyapp.fluently.common.model

import kotlinx.serialization.Serializable

@Serializable
data class Lesson(
    val wordsPerLesson: Int,
    val lessonId: String,
    val components: List<LessonComponent>,
    val currentLessonComponentIndex: Int = 0
) {
    val currentComponent: LessonComponent
        get() = components[currentLessonComponentIndex]
}