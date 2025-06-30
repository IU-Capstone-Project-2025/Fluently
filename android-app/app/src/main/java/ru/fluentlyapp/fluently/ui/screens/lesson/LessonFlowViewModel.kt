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
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.app.navigation.Destination
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordObserver
import javax.inject.Inject

@HiltViewModel
class LessonFlowViewModel @Inject constructor(
    savedStateHandle: SavedStateHandle,
    private val lessonRepository: LessonRepository
) : ViewModel() {
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
                    val fetchedLesson =
                        fetchLesson(lessonId) // the potential exception will silently cancel the coroutine
                    saveLesson(fetchedLesson)
                }
            }
        }
    }

    private inline fun <reified T> checkCurrentComponentOrNull(): T? {
        return lesson.value?.currentComponent as? T
    }

    fun moveToNextComponent() {
        // Save the current value of lesson and the component to variables
        val currentLesson = lesson.value
        val currentComponent = currentLesson?.currentComponent

        if (currentLesson == null || currentComponent == null)
            return

        if (currentComponent !is Exercise || currentComponent.isAnswered) {
            // Copy the current lesson and just increase the lesson component index
            val updatedLesson = currentLesson.copy(
                // Additional check to not go beyond the list size
                currentLessonComponentIndex = (currentLesson.currentLessonComponentIndex + 1).coerceIn(
                    0,
                    currentLesson.components.size - 1
                )
            )

            // Save the updated lesson
            viewModelScope.launch {
                lessonRepository.saveLesson(updatedLesson)
            }
        }
    }

    // Observer that is responsible for handling "choose translation" exercises
    val chooseTranslationObserver = object : ChooseTranslationObserver() {
        override fun onVariantClick(variantIndex: Int) {
            val currentComponent =
                checkCurrentComponentOrNull<Exercise.ChooseTranslation>() ?: return
            val updatedComponent = currentComponent.copy(selectedVariant = variantIndex)
            updateCurrentComponent(updatedComponent)
        }

        override fun onCompleteExercise() = moveToNextComponent()
    }

    // Observer for the "learn new word" exercises
    val newWordObserver = object : NewWordObserver() {
        override fun onUserKnowsWord(doesUserKnowWord: Boolean) {
            val currentComponent = checkCurrentComponentOrNull<Exercise.NewWord>() ?: return
            val updatedComponent = currentComponent.copy(doesUserKnow = doesUserKnowWord)
            updateCurrentComponent(updatedComponent)
        }

        override fun onCompleteExercise() = moveToNextComponent()
    }

    // Observer for fill the gap exercises
    val fillGapsObserver = object : FillGapsObserver() {
        override fun onCompleteExercise() = moveToNextComponent()

        override fun onVariantClick(variantIndex: Int) {
            val currentComponent = checkCurrentComponentOrNull<Exercise.FillTheGap>() ?: return
            val updatedComponent = currentComponent.copy(selectedVariant = variantIndex)
            updateCurrentComponent(updatedComponent)
        }
    }

    // Observer for input word exercise
    val inputWordObserver = object : InputWordObserver {
        override fun onCompleteExercise() = moveToNextComponent()

        override fun onConfirmInput(input: String) {
            val currentComponent = checkCurrentComponentOrNull<Exercise.InputWord>() ?: return
            val updatedComponent = currentComponent.copy(inputtedWord = input)
            updateCurrentComponent(updatedComponent)
        }
    }

    private fun updateCurrentComponent(newComponent: LessonComponent) {
        // Update the component at the currentComponentIndex
        val currentLesson = lesson.value
        if (currentLesson == null)
            return

        // Copy the current lesson's prefix and suffix and inject the `newComponent`
        val newComponents = with(currentLesson) {
            buildList {
                addAll(components.subList(0, currentLessonComponentIndex))
                add(newComponent)
                addAll(components.subList(currentLessonComponentIndex + 1, components.size))
            }
        }
        val newLesson = currentLesson.copy(components = newComponents)

        viewModelScope.launch {
            lessonRepository.saveLesson(newLesson)
        }
    }
}