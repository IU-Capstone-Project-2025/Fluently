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

@HiltViewModel
class HomeScreenViewModel @Inject constructor(
    private val lessonRepository: LessonRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(HomeScreenUiState())
    val uiState = _uiState.asStateFlow()

    init {
        viewModelScope.launch {
            lessonRepository.getSavedOngoingLessonIdAsFlow().collect { id ->
                if (id == null) {
                    _uiState.update {
                        it.copy(ongoingLessonState = OngoingLessonState.NOT_STARTED)
                    }
                }
            }
        }
    }

    fun fetchCurrentLessonForUser(
        onSuccessfulFetch: (lessonId: String) -> Unit
    ) {
        _uiState.update {
            it.copy(ongoingLessonState = OngoingLessonState.LOADING)
        }
        viewModelScope.launch(Dispatchers.IO) {
            val lessonId = lessonRepository.getSavedOngoingLessonIdAsFlow().first()

            if (lessonId != null) {
                onSuccessfulFetch(lessonId)
            }

            try {
                val ongoingLesson = lessonRepository.fetchCurrentLesson()
                lessonRepository.setSavedOngoingLessonId(ongoingLesson.lessonId)

                onSuccessfulFetch(ongoingLesson.lessonId)
            } catch (ex: Exception) {
                Log.e("HomeScreenViewModel", "Failed to catch lesson: $ex")
                _uiState.update {
                    it.copy(ongoingLessonState = OngoingLessonState.ERROR)
                }

            }
        }
    }
}