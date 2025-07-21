package ru.fluentlyapp.fluently.ui.screens.onboarding

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.receiveAsFlow
import kotlinx.coroutines.flow.update
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.feature.topics.TopicRepository
import ru.fluentlyapp.fluently.feature.userpreferences.UserPreferencesRepository
import ru.fluentlyapp.fluently.utils.safeLaunch
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class OnboardingViewModel @Inject constructor(
    private val userPreferencesRepository: UserPreferencesRepository,
    private val topicRepository: TopicRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow<OnboardingScreenUiState>(
        OnboardingScreenUiState(
            initialLoadingState = InitialLoadingState.LOADING,
            userPreferences = UserPreferences.empty(),
            availableTopics = emptyList(),
            uploadingLoadingState = UploadingLoadingState.IDLE
        )
    )
    val uiState = _uiState.asStateFlow()

    private val _commands = Channel<OnboardingScreenCommand>()
    val commands = _commands.receiveAsFlow()

    init {
        initOnboarding()

        viewModelScope.safeLaunch {
            userPreferencesRepository.getCachedUserPreferences().collect { preferences ->
                if (preferences != null) {
                    _uiState.update {
                        it.copy(userPreferences = preferences)
                    }
                }
            }
        }
    }

    fun initOnboarding() {
        viewModelScope.safeLaunch {
            _uiState.update {
                it.copy(initialLoadingState = InitialLoadingState.LOADING)
            }
            try {
                val remotePreferences = userPreferencesRepository.getRemoteUserPreferences()
                val topics = topicRepository.getAvailableTopic()
                _uiState.update {
                    it.copy(
                        userPreferences = remotePreferences,
                        availableTopics = topics,
                        initialLoadingState = InitialLoadingState.SUCCESS
                    )
                }
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update {
                    it.copy(
                        initialLoadingState = InitialLoadingState.ERROR
                    )
                }
            }
        }
    }

    fun updateUserPreferences(userPreferences: UserPreferences) {
        viewModelScope.safeLaunch {
            userPreferencesRepository.updateCachedUserPreferences(userPreferences)
        }
    }

    fun completeUserPreferences() {
        _uiState.update {
            it.copy(uploadingLoadingState = UploadingLoadingState.UPLOADING)
        }
        viewModelScope.safeLaunch {
            try {
                _uiState.value.userPreferences.let {
                    userPreferencesRepository.updateRemoteUserPreferences(it)
                    userPreferencesRepository.updateCachedUserPreferences(it)
                }
                _uiState.update {
                    it.copy(uploadingLoadingState = UploadingLoadingState.SUCCESS)
                }
                _commands.send(OnboardingScreenCommand.UserPreferencesUploadedCommand)
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update {
                    it.copy(uploadingLoadingState = UploadingLoadingState.ERROR)
                }
            }
        }
    }
}