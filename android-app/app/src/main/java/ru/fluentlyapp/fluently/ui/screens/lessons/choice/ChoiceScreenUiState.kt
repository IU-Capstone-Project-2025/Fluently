package ru.fluentlyapp.fluently.ui.screens.lessons.choice

data class ChoiceScreenUiState(
    val word: String,
    val variants: List<String>,
    val correctVariant: Int,
    val selectedVariant: Int? = null
)