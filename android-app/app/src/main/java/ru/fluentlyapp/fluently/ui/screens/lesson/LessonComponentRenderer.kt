package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.animation.togetherWith
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.other.LoadingLessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.FillGapsObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.InputWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordObserver
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordExercise

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
    inputWordObserver: InputWordObserver
) {
    AnimatedContent(
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
                LoadingLessonComponent(modifier = modifier)
            }

            is Exercise.ChooseTranslation -> {
                ChooseTranslationExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    chooseTranslationObserver = chooseTranslationObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.NewWord -> {
                NewWordExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    newWordObserver = newWordObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.FillTheGap -> {
                FillGapsExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    fillGapsObserver = fillGapsObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.InputWord -> {
                InputWordExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    observer = inputWordObserver,
                    isCompleted = targetComponent.isAnswered
                )
            }
        }
    }
}