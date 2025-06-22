package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.animation.AnimatedContent
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.other.LoadingLessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationController
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.ChooseTranslationExercise
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordController
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.NewWordExercise

@Composable
fun LessonComponentRenderer(
    modifier: Modifier,
    component: LessonComponent,
    chooseTranslationController: ChooseTranslationController,
    newWordController: NewWordController,
) {
    when (component) {
        is LessonComponent.Loading -> {
            LoadingLessonComponent(modifier = modifier)
        }

        is Exercise.ChooseTranslation -> {
            ChooseTranslationExercise(
                modifier = modifier,
                exerciseState = component,
                chooseTranslationController = chooseTranslationController,
                isCompleted = component.isAnswered
            )
        }

        is Exercise.NewWord -> {
            NewWordExercise(
                modifier = modifier,
                exerciseState = component,
                newWordController = newWordController,
                isCompleted = component.isAnswered
            )
        }
    }
}