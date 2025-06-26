package ru.fluentlyapp.fluently.common.model

data class Lesson(
    val lessonId: String,
    val components: List<LessonComponent>,
    val currentLessonComponentIndex: Int = 0
) {
    val currentComponent: LessonComponent
        get() = components[currentLessonComponentIndex]
}