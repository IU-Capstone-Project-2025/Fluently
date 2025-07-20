package ru.fluentlyapp.fluently.ui.screens.settings

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.auth.AuthManager
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.feature.topics.TopicRepository
import ru.fluentlyapp.fluently.feature.userpreferences.UserPreferencesRepository
import ru.fluentlyapp.fluently.utils.safeLaunch
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class SettingsScreenViewModel @Inject constructor(
    private val authManager: AuthManager,
    private val userPreferencesRepository: UserPreferencesRepository,
    private val topicRepository: TopicRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow<SettingsScreenUiState>(
        SettingsScreenUiState(
            userPreferences = UserPreferences.empty(),
            uploadingState = SettingsUploading.IDLE,
            availableTopics = emptyList()
        )
    )
    val uiState = _uiState.asStateFlow()

    private val _commands = Channel<SettingScreenCommand>()
    val commands = _commands.receiveAsFlow()

    init {
        viewModelScope.launch {
            safeLaunch {
                val remotePreferences = userPreferencesRepository.getRemoteUserPreferences()
                userPreferencesRepository.updateCachedUserPreferences(remotePreferences)
            }

            safeLaunch {
                _uiState.update {
                    it.copy(
                        availableTopics = topicRepository.getAvailableTopic()
                    )
                }
            }

            safeLaunch {
                userPreferencesRepository.getCachedUserPreferences().collect { preferences ->
                    if (preferences != null) {
                        _uiState.update {
                            it.copy(userPreferences = preferences)
                        }
                    }
                }
            }

        }
    }

    fun logout() {
        viewModelScope.safeLaunch {
            authManager.deleteServerToken()
            _commands.send(SettingScreenCommand.LoginCredentialsRemovedCommand)
        }
    }

    fun updateUserPreferences(preferences: UserPreferences) {
        _uiState.update {
            it.copy(
                userPreferences = preferences
            )
        }
    }

    fun uploadUserPreferences() {
        _uiState.update {
            it.copy(uploadingState = SettingsUploading.UPLOADING)
        }
        viewModelScope.safeLaunch {
            try {
                _uiState.value.userPreferences.also {
                    userPreferencesRepository.updateRemoteUserPreferences(it)
                    userPreferencesRepository.updateCachedUserPreferences(it)
                }
                _uiState.update {
                    it.copy(uploadingState = SettingsUploading.SUCCESS)
                }
                _commands.send(SettingScreenCommand.SettingsUpdatedCommand)
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update {
                    it.copy(uploadingState = SettingsUploading.ERROR)
                }
            }
        }
    }
}