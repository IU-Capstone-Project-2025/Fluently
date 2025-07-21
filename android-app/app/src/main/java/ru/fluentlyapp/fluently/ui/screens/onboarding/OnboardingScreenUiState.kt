package ru.fluentlyapp.fluently.ui.screens.onboarding

import ru.fluentlyapp.fluently.common.model.UserPreferences

enum class InitialLoadingState {
    LOADING,
    SUCCESS,
    ERROR
}

enum class UploadingLoadingState {
    IDLE,
    SUCCESS,
    ERROR,
    UPLOADING
}

data class OnboardingScreenUiState(
    val initialLoadingState: InitialLoadingState,
    val uploadingLoadingState: UploadingLoadingState,
    val userPreferences: UserPreferences,
    val availableTopics: List<String>
)