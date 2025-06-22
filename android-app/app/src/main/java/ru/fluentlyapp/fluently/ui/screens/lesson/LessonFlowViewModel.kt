package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.lifecycle.SavedStateHandle
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import androidx.navigation.toRoute
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.navigation.Destination
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationController
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordController
import javax.inject.Inject

@HiltViewModel
class LessonFlowViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val lessonRepository: LessonRepository
): ViewModel() {
    // Get the lessonId of the lesson this flow represents
    private val lessonId = savedStateHandle.toRoute<Destination.LessonScreen>().lessonId

    // Start collecting the lesson state from the repository
    val lesson = lessonRepository.getSavedLessonAsFlow(lessonId)
        .stateIn(
            viewModelScope,
            SharingStarted.Eagerly,
            initialValue = null
        )

    init {
        viewModelScope.launch {
            // There is possibility that the lesson with lessonId hasn't been fetched yet
            with(lessonRepository) {
                if (getSavedLesson(lessonId) == null) { // the lesson hasn't been fetched yet
                    val fetchedLesson = fetchLesson(lessonId) // the potential exception will silently cancel the coroutine
                    saveLesson(fetchedLesson)
                }
            }
        }
    }

    private inline fun <reified T> checkCurrentComponentOrNull(): T? {
        return lesson.value?.currentComponent as? T
    }

    fun moveToNextComponent() {
        val currentComponent = lesson.value?.currentComponent
        if (currentComponent == null) return
        if (currentComponent !is Exercise || currentComponent.isAnswered) {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent(lessonId)
            }
        }
    }

    // Controller that is responsible for handling "choose translation" exercises
    val chooseTranslationController = object : ChooseTranslationController() {
        override fun onVariantClick(variantIndex: Int) {
            val currentComponent = checkCurrentComponentOrNull<Exercise.ChooseTranslation>() ?: return
            val updatedComponent = currentComponent.copy(selectedVariant = variantIndex)
            updateCurrentComponent(updatedComponent)
        }

        override fun onCompleteExercise() = moveToNextComponent()
    }

    // Controller for the "learn new word" exercises
    val newWordController = object : NewWordController() {
        override fun onUserKnowsWord(doesUserKnowWord: Boolean) {
            val currentComponent = checkCurrentComponentOrNull<Exercise.NewWord>() ?: return
            val updatedComponent = currentComponent.copy(doesUserKnow = doesUserKnowWord)
            updateCurrentComponent(updatedComponent)
        }

        override fun onCompleteExercise() = moveToNextComponent()
    }

    private fun updateCurrentComponent(newComponent: LessonComponent) {
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