package ru.fluentlyapp.fluently.ui.screens.wordsprogress

import ru.fluentlyapp.fluently.ui.theme.components.WordUiState

data class WordsProgressUiState(
    val words: List<WordUiState>,
    val pageTitle: String,
    val searchString: String
)