package ru.fluentlyapp.fluently.ui.screens.lesson

import android.util.Log
import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.animation.togetherWith
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.Modifier
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Dialog
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.FinishDecoration
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.FinishDecorationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.LoadingDecoration
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.OnboardingDecorationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.OnboardingDecorationUiState
import ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration.OnboardingDecoration
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.DialogExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.DialogObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordExercise
import ru.fluentlyapp.fluently.ui.screens.login.LoginScreen
import timber.log.Timber

data class LessonComponentWithIndex(
    val lessonComponent: LessonComponent,
    val index: Int,
)

@Composable
fun LessonComponentRenderer(
    modifier: Modifier,
    component: LessonComponentWithIndex,
    chooseTranslationObserver: ChooseTranslationObserver,
    newWordObserver: NewWordObserver,
    fillGapsObserver: FillGapsObserver,
    inputWordObserver: InputWordObserver,
    onboardingDecorationObserver: OnboardingDecorationObserver,
    finishDecorationObserver: FinishDecorationObserver,
    dialogObserver: DialogObserver
) {
    AnimatedContent(
        modifier = modifier,
        targetState = component,
        transitionSpec = {
            slideInHorizontally(
                tween(500),
                initialOffsetX = { it }
            ) togetherWith slideOutHorizontally(
                tween(500),
                targetOffsetX = { -it }
            )
        },
        contentKey = { it.index }
    ) { (targetComponent, index) ->
        when (targetComponent) {
            is Decoration.Loading -> {
                LoadingDecoration(modifier = Modifier.fillMaxSize())
            }

            is Decoration.Finish -> {
                FinishDecoration(
                    modifier = Modifier.fillMaxSize(),
                    finishDecorationObserver = finishDecorationObserver
                )
            }

            is Decoration.Onboarding -> {
                OnboardingDecoration(
                    modifier = Modifier.fillMaxSize(),
                    onboardingDecorationUiState = OnboardingDecorationUiState(
                        newWordsCount = targetComponent.wordsToBeLearned,
                        exercisesCount = targetComponent.featuredExercises
                    ),
                    onboardingDecorationObserver = onboardingDecorationObserver
                )
            }

            is Exercise.ChooseTranslation -> {
                ChooseTranslationExercise(
                    modifier = Modifier.fillMaxSize(),
                    exerciseState = targetComponent,
                    chooseTranslationObserver = chooseTranslationObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.NewWord -> {
                NewWordExercise(
                    modifier = Modifier.fillMaxSize(),
                    exerciseState = targetComponent,
                    newWordObserver = newWordObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.FillTheGap -> {
                FillGapsExercise(
                    modifier = Modifier.fillMaxSize(),
                    exerciseState = targetComponent,
                    fillGapsObserver = fillGapsObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.InputWord -> {
                InputWordExercise(
                    modifier = Modifier.fillMaxSize(),
                    exerciseState = targetComponent,
                    observer = inputWordObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Dialog -> {
                DialogExercise(
                    modifier = Modifier.fillMaxSize(),
                    exerciseState = targetComponent,
                    dialogObserver = dialogObserver,
                    isCompleted = targetComponent.isFinished
                )
            }
        }
    }
}