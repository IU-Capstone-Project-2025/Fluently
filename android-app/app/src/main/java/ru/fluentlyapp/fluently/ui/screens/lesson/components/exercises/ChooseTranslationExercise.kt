package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
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
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.ui.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.FluentlyColors
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews
import ru.fluentlyapp.fluently.ui.utils.MediumPhonePreview
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview

abstract class ChooseTranslationObserver {
    abstract fun onVariantClick(variantIndex: Int)
    abstract fun onCompleteExercise()
}

@Composable
fun ChooseTranslationExercise(
    modifier: Modifier = Modifier,
    exerciseState: Exercise.ChooseTranslation,
    chooseTranslationObserver: ChooseTranslationObserver,
    isCompleted: Boolean
) {
    Column(
        modifier = modifier.background(FluentlyTheme.colors.surface),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Column(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
                .verticalScroll(state = rememberScrollState())
                .padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(exerciseState.word, fontSize = 32.sp, textAlign = TextAlign.Center)

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                stringResource(R.string.choose_correct_translation),
                fontSize = 16.sp
            )

            Spacer(modifier = Modifier.height(8.dp))

            val correctColor = FluentlyTheme.colors.correct
            val correctAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = correctColor,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val wrongColor = FluentlyTheme.colors.wrong
            val wrongAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = wrongColor,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val notChosenModifier = remember {
                Modifier.alpha(alpha = .4f)
            }

            repeat(times = exerciseState.answerVariants.size) { index ->
                val itemModifier = when {
                    exerciseState.selectedVariant == null -> Modifier.clickable(
                        onClick = { chooseTranslationObserver.onVariantClick(index) },
                    )

                    exerciseState.correctVariant == index -> correctAnswerModifier
                    exerciseState.selectedVariant == index -> wrongAnswerModifier
                    else -> notChosenModifier
                }

                Text(
                    modifier = Modifier then itemModifier
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


        Box(
            modifier = Modifier.fillMaxWidth().height(160.dp),
            contentAlignment = Alignment.Center
        ) {
            if (isCompleted) {
                ExerciseContinueButton(
                    onClick = { chooseTranslationObserver.onCompleteExercise() }
                )
            }
        }
    }
}

@SmallPhonePreview
@MediumPhonePreview
@Composable
fun ChooseTranslationExercisePreview() {
    FluentlyTheme {
        ChooseTranslationExercise(
            modifier = Modifier
                .fillMaxSize(),
            exerciseState = Exercise.ChooseTranslation(
                word = "Pretty long word aba aba wow this is very long",
                answerVariants = listOf(
                    "Влияние",
                    "Очень длинный вариант ответа капец он реально очень длинный господи ну почему он такой длинный",
                    "Двойственность",
                    "Комар"
                ),
                correctVariant = 0,
                selectedVariant = 2,
            ),
            chooseTranslationObserver = object : ChooseTranslationObserver() {
                override fun onVariantClick(variantIndex: Int) { }
                override fun onCompleteExercise() { }
            },
            isCompleted = true
        )
    }
}