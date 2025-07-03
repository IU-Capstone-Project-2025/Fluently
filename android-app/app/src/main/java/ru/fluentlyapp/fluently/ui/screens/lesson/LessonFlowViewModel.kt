package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordObserver
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class LessonFlowViewModel @Inject constructor(
    private val lessonRepository: LessonRepository
) : ViewModel() {
    // Start collecting the lesson state from the repository
    val currentComponent: StateFlow<LessonComponent?> = lessonRepository.currentComponent()
        .stateIn(
            viewModelScope,
            SharingStarted.Eagerly,
            initialValue = null
        )

    /**
     * Check if the current component is an instance of type `T`, and if it is, then try
     * to save the result of the block lambda into the `lessonRepository`.
     */
    private inline fun <reified T : LessonComponent> safeApplyAndUpdate(
        produceNewComponent: (oldComponent: T) -> T
    ) {
        (currentComponent.value as? T)?.let { oldComponent ->
            produceNewComponent(oldComponent).also { newComponent ->
                try {
                    viewModelScope.launch {
                        lessonRepository.updateCurrentComponent(newComponent)
                    }
                } catch (ex: Exception) {
                    Timber.e(ex)
                }
            }
        }
    }

    // Observer that is responsible for handling the "choose translation" exercises
    val chooseTranslationObserver = object : ChooseTranslationObserver() {
        override fun onVariantClick(variantIndex: Int) {
            safeApplyAndUpdate<Exercise.ChooseTranslation> {
                it.copy(selectedVariant = variantIndex)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for the "learn new word" exercises
    val newWordObserver = object : NewWordObserver() {
        override fun onUserKnowsWord(doesUserKnowWord: Boolean) {
            safeApplyAndUpdate<Exercise.NewWord> {
                it.copy(doesUserKnow = doesUserKnowWord)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for fill the gap exercises
    val fillGapsObserver = object : FillGapsObserver() {

        override fun onVariantClick(variantIndex: Int) {
            safeApplyAndUpdate<Exercise.FillTheGap> {
                it.copy(selectedVariant = variantIndex)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }

    // Observer for input word exercise
    val inputWordObserver = object : InputWordObserver {
        override fun onConfirmInput(inputtedWord: String) {
            safeApplyAndUpdate<Exercise.InputWord> {
                it.copy(inputtedWord = inputtedWord)
            }
        }

        override fun onCompleteExercise() {
            viewModelScope.launch {
                lessonRepository.moveToNextComponent()
            }
        }
    }
}