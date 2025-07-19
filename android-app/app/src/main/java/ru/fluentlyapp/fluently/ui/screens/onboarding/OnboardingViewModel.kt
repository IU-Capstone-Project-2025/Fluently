package ru.fluentlyapp.fluently.ui.screens.onboarding

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.common.model.CefrLevel
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
            loadingState = LoadingState.LOADING,
            userPreferences = null,
            availableTopics = emptyList()
        )
    )
    val uiState = _uiState.asStateFlow()

    init {
        initOnboarding()
    }

    fun initOnboarding() {
        viewModelScope.safeLaunch {
            _uiState.update {
                it.copy(loadingState = LoadingState.LOADING)
            }
            try {
                val remotePreferences = userPreferencesRepository.getRemoteUserPreferences()
                val topics = topicRepository.getAvailableTopic()
                _uiState.update {
                    it.copy(
                        userPreferences = remotePreferences,
                        availableTopics = topics,
                        loadingState = LoadingState.SUCCESS
                    )
                }
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update {
                    it.copy(
                        loadingState = LoadingState.ERROR
                    )
                }
            }
        }
    }

    fun updateUserPreferences(userPreferences: UserPreferences) {
        _uiState.update { it.copy(userPreferences = userPreferences) }
    }

    fun sendUserPreferences(preferences: UserPreferences) {
        viewModelScope.safeLaunch {
            userPreferencesRepository.updateRemoteUserPreferences(preferences)
        }
    }
}