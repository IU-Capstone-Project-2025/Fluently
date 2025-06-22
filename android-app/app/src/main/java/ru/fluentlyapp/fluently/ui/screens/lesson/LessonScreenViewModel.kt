package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.toRoute
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.conflate
import kotlinx.coroutines.flow.onEach
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.Lesson
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.navigation.Destination
import javax.inject.Inject

@HiltViewModel
class LessonScreenViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val lessonRepository: LessonRepository
) : ViewModel() {
    private val _uiState = MutableStateFlow(
        LessonScreenUiState(
            currentComponent = LessonComponent.Loading,
            showContinueButton = false
        )
    )
    val uiState = _uiState.asStateFlow()

    private val lessonId = savedStateHandle.toRoute<Destination.LessonScreen>().lessonId
    private val lesson = lessonRepository.getSavedLessonAsFlow(lessonId)
        .conflate()
        .onEach {
            if (it != null) parseLesson(it)
        }
        .stateIn(
            viewModelScope,
            SharingStarted.Eagerly,
            initialValue = null
        )

    init {
        viewModelScope.launch {
            with(lessonRepository) {
                if (getSavedLesson(lessonId) == null) {
                    saveLesson(fetchLesson(lessonId))
                }
            }
        }
    }

    /**
     * For lesson, find the current component and put it into the uiState
     */
    fun parseLesson(lesson: Lesson) {
        _uiState.update {
            val currentComponent = lesson.components[lesson.currentLessonComponentIndex]
            it.copy(
                currentComponent = currentComponent,
                showContinueButton = currentComponent !is Exercise || currentComponent.isAnswered
            )
        }
    }

    fun moveToNextComponent() {
        viewModelScope.launch {
            lessonRepository.moveToNextComponent(lessonId)
        }
    }

    fun updateCurrentComponent(newComponent: LessonComponent) {
        val previousLesson = lesson.value
        if (previousLesson == null)
            return

        val newComponents = with(previousLesson) {
            buildList {
                addAll(components.subList(0, currentLessonComponentIndex))
                add(newComponent)
                addAll(components.subList(currentLessonComponentIndex + 1, components.size))
            }
        }
        val newLesson = previousLesson.copy(components = newComponents)

        viewModelScope.launch {
            lessonRepository.saveLesson(newLesson)
        }
    }
}
