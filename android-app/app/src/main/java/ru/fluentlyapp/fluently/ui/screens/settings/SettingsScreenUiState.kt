package ru.fluentlyapp.fluently.ui.screens.settings

import ru.fluentlyapp.fluently.common.model.UserPreferences

enum class SettingsUploading {
    IDLE,
    UPLOADING,
    ERROR,
    SUCCESS
}

data class SettingsScreenUiState(
    val userPreferences: UserPreferences,
    val uploadingState: SettingsUploading,
    val availableTopics: List<String>
)