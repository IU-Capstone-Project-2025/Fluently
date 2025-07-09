package ru.fluentlyapp.fluently.ui.screens.home

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import javax.inject.Inject
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.receiveAsFlow
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreenUiState.OngoingLessonState
import timber.log.Timber


@HiltViewModel
class HomeScreenViewModel @Inject constructor(
    private val lessonRepository: LessonRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(HomeScreenUiState())
    val uiState = _uiState.asStateFlow()

    private val _commandsChannel = Channel<HomeCommands>()
    val commandsChannel = _commandsChannel.receiveAsFlow()

    init {
        viewModelScope.launch {
            val currentComponent = lessonRepository.currentComponent().first()
            if (currentComponent != null) {
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
            }
        }
    }

    fun ensureOngoingLesson() {
        _uiState.update {
            it.copy(ongoingLessonState = OngoingLessonState.LOADING)
        }

        viewModelScope.launch {
            try {
                if (lessonRepository.hasSavedLesson()) {
                    _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
                    _commandsChannel.send(HomeCommands.NavigateToLesson)
                    return@launch
                }

                lessonRepository.fetchAndSaveOngoingLesson()
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
                _commandsChannel.send(HomeCommands.NavigateToLesson)
            } catch (ex: Exception) {
                Timber.e(ex)
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.ERROR) }
            }
        }
    }
}