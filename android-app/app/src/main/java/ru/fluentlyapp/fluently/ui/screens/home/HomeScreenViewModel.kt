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

@HiltViewModel
class HomeScreenViewModel @Inject constructor(
    private val lessonRepository: LessonRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(HomeScreenUiState())
    val uiState = _uiState.asStateFlow()

    /**
     * When the user presses on the new lesson, the received current lesson id will be emitted
     * into this flow.
     */
    private val _ongoingLessonId = MutableStateFlow<String?>(null)
    val ongoingLessonId = _ongoingLessonId.asStateFlow()

    init {
        viewModelScope.launch {
            // On each update of the saved lesson id, update the ui state correspondingly
            lessonRepository.getSavedOngoingLessonIdAsFlow().collect { id ->
                if (id == null) {
                    // No saved id -> mark the absence of any saved ongoing lesson
                    _uiState.update {
                        it.copy(ongoingLessonState = OngoingLessonState.NOT_STARTED)
                    }
                } else {
                    _uiState.update {
                        it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED)
                    }
                }
            }
        }
    }

    fun getOngoingLessonId() {
        _uiState.update {
            it.copy(ongoingLessonState = OngoingLessonState.LOADING)
        }

        viewModelScope.launch(Dispatchers.IO) {
            val lessonId = lessonRepository.getSavedOngoingLessonIdAsFlow().first()

            if (lessonId != null) {
                _ongoingLessonId.value = lessonId
                _uiState.update {
                    it.copy(ongoingLessonState = OngoingLessonState.HAS_PAUSED)
                }
                return@launch
            }

            try {
                val ongoingLesson = lessonRepository.fetchCurrentLesson()
                lessonRepository.setSavedOngoingLessonId(ongoingLesson.lessonId)
                _ongoingLessonId.value = ongoingLesson.lessonId
            } catch (ex: Exception) {
                Log.e("HomeScreenViewModel", "Failed to catch lesson: $ex")
                _uiState.update {
                    it.copy(ongoingLessonState = OngoingLessonState.ERROR)
                }
            }
        }
    }
}