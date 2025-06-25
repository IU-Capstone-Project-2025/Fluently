package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.animation.togetherWith
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.other.LoadingLessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationController
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordController
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordExercise

data class LessonComponentWithIndex(
    val lessonComponent: LessonComponent,
    val index: Int,
)

@Composable
fun LessonComponentRenderer(
    modifier: Modifier,
    component: LessonComponentWithIndex,
    chooseTranslationController: ChooseTranslationController,
    newWordController: NewWordController,
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
            is LessonComponent.Loading -> {
                LoadingLessonComponent(modifier = modifier)
            }

            is Exercise.ChooseTranslation -> {
                ChooseTranslationExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    chooseTranslationController = chooseTranslationController,
                    isCompleted = targetComponent.isAnswered
                )
            }

            is Exercise.NewWord -> {
                NewWordExercise(
                    modifier = modifier,
                    exerciseState = targetComponent,
                    newWordController = newWordController,
                    isCompleted = targetComponent.isAnswered
                )
            }
        }
    }
}