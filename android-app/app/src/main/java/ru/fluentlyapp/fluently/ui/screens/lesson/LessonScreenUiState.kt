package ru.fluentlyapp.fluently.ui.screens.lesson

import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent

data class LessonScreenUiState(
    val currentComponent: LessonComponent,
    val showContinueButton: Boolean
)