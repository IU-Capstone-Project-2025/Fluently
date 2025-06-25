package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.ui.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

abstract class ChooseTranslationController {
    abstract fun onVariantClick(variantIndex: Int)
    abstract fun onCompleteExercise()
}

@Composable
fun ChooseTranslationExercise(
    modifier: Modifier = Modifier,
    exerciseState: Exercise.ChooseTranslation,
    chooseTranslationController: ChooseTranslationController,
    isCompleted: Boolean
) {
    Column(
        modifier = modifier,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Column(
            modifier = Modifier
                .background(color = FluentlyTheme.colors.surface)
                .weight(1f)
                .fillMaxWidth()
                .verticalScroll(
                    state = rememberScrollState()
                ),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(exerciseState.word, fontSize = 32.sp)

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                stringResource(R.string.choose_correct_translation),
                fontSize = 16.sp
            )

            Spacer(modifier = Modifier.height(8.dp))

            val correctAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = Color.Green,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val wrongAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = Color.Red,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val notChosenModifier = remember {
                Modifier.alpha(alpha = .4f)
            }

            repeat(times = exerciseState.answerVariants.size) { index ->
                val itemModifier = when {
                    exerciseState.selectedVariant == null -> Modifier.clickable(
                        onClick = { chooseTranslationController.onVariantClick(index) },
                    )

                    exerciseState.correctVariant == index -> correctAnswerModifier
                    exerciseState.selectedVariant == index -> wrongAnswerModifier
                    else -> notChosenModifier
                }

                Text(
                    modifier = Modifier
                            then itemModifier
                        .clip(shape = RoundedCornerShape(12.dp))
                        .fillMaxWidth(fraction = 0.8f)
                        .background(color = FluentlyTheme.colors.surfaceContainerHigh)
                        .padding(16.dp),
                    text = exerciseState.answerVariants[index],
                )
                if (index != exerciseState.answerVariants.size - 1) {
                    Spacer(modifier = Modifier.height(4.dp))
                }
            }
        }

        ExerciseContinueButton(
            modifier = Modifier.fillMaxWidth(.7f).padding(32.dp).height(80.dp),
            enabled = isCompleted,
            onClick = chooseTranslationController::onCompleteExercise
        )
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun ChooseTranslationExercisePreview() {
    FluentlyTheme {
        ChooseTranslationExercise(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface),
            exerciseState = Exercise.ChooseTranslation(
                word = "Influence",
                answerVariants = listOf(
                    "Влияние", "Благодарность", "Двойственность", "Комар"
                ),
                correctVariant = 0,
                selectedVariant = 2,
            ),
            chooseTranslationController = object : ChooseTranslationController() {
                override fun onVariantClick(variantIndex: Int) { }
                override fun onCompleteExercise() { }
            },
            isCompleted = true
        )
    }
}