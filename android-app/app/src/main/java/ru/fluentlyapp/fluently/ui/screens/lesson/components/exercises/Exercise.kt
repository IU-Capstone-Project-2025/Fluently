package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent

@Composable
fun Exercise(
    modifier: Modifier,
    exercise: Exercise,
    onUpdateExercise: (Exercise) -> Unit
) {
    if (exercise is Exercise.NewWord) {
        NewWordExercise(
            modifier = modifier,
            exerciseState = exercise,
            onLearnWordClick = {
                onUpdateExercise(
                    exercise.copy(doesUserKnow = false)
                )
            },
            onKnowWordClick = {
                onUpdateExercise(
                    exercise.copy(
                        doesUserKnow = true
                    )
                )
            }
        )
    } else if (exercise is Exercise.ChooseTranslation) {
        ChooseTranslationExercise(
            modifier = modifier,
            exerciseState = exercise,
            onVariantClick = { selectedVariant ->
                onUpdateExercise(
                    exercise.copy(selectedVariant = selectedVariant)
                )
            }
        )
    }
}