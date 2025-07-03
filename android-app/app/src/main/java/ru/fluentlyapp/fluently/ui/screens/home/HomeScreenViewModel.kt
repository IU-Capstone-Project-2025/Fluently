package ru.fluentlyapp.fluently.ui.screens.home

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import javax.inject.Inject
import kotlinx.coroutines.flow.first
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreenUiState.OngoingLessonState
import timber.log.Timber

@HiltViewModel
class HomeScreenViewModel @Inject constructor(
    private val lessonRepository: LessonRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(HomeScreenUiState())
    val uiState = _uiState.asStateFlow()

    /**
     * When the user presses on the new lesson, this flag will indicate whether the ongoing lesson
     * is prepared.
     */
    private val _ongoingLessonIsReady = MutableStateFlow<Boolean>(false)
    val ongoingLessonIsReady = _ongoingLessonIsReady.asStateFlow()

    init {
        viewModelScope.launch {
            val currentComponent = lessonRepository.currentComponent().first()
            if (currentComponent != null) {
                _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED) }
            }
        }
    }

    /**
     * If the ongoing lesson is not loaded, load it from the server. If the loading is
     * successful or the lesson is already stored, emit true into the `_ongoingLessonIsReady` flow.
     */
    fun ensureOngoingLesson() {
        _uiState.update {
            it.copy(ongoingLessonState = OngoingLessonState.LOADING)
        }

        try {
            viewModelScope.launch {
                if (lessonRepository.hasSavedLesson()) {
                    _ongoingLessonIsReady.value = true
                    return@launch
                }
                lessonRepository.fetchAndSaveOngoingLesson()
                _ongoingLessonIsReady.value = true
            }
        } catch (ex: Exception) {
            Timber.e(ex)
            _uiState.update { it.copy(ongoingLessonState = OngoingLessonState.ERROR) }
        }
    }
}