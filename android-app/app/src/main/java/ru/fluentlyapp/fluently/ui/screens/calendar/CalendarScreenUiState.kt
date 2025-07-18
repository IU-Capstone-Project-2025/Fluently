package ru.fluentlyapp.fluently.ui.screens.calendar

import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import java.time.LocalDate

data class CalendarScreenUiState(
    val day: LocalDate,
    val learnedWords: List<WordUiState>,
    val inProgressWords: List<WordUiState>,
    val showIsLearning: Boolean
)