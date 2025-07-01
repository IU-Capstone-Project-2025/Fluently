package ru.fluentlyapp.fluently.common.model

interface LessonComponent {
    object Loading : LessonComponent
    object Finished : LessonComponent
}