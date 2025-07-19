package ru.fluentlyapp.fluently.ui.screens.onboarding

import ru.fluentlyapp.fluently.common.model.UserPreferences

enum class LoadingState {
    LOADING,
    SUCCESS,
    ERROR
}

data class OnboardingScreenUiState(
    val loadingState: LoadingState,
    val userPreferences: UserPreferences?,
    val availableTopics: List<String>
)