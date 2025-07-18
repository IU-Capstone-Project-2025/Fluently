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
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.style.TextDecoration
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.ui.theme.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

abstract class FillGapsObserver {
    abstract fun onVariantClick(variantIndex: Int)
    abstract fun onCompleteExercise()
}

@Composable
fun FillGapsExercise(
    modifier: Modifier = Modifier,
    exerciseState: Exercise.FillTheGap,
    fillGapsObserver: FillGapsObserver,
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
                .verticalScroll(
                    state = rememberScrollState()
                )
                .padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            if (exerciseState.selectedVariant == null) {
                Text(
                    text = exerciseState.sentence.joinToString(" ____ "),
                    fontSize = 20.sp
                )
            } else {
                Text(
                    text = run {
                        val missedWord = exerciseState.answerVariants[exerciseState.correctVariant]
                        buildAnnotatedString {
                            for ((index, part) in exerciseState.sentence.withIndex()) {
                                append(part)
                                if (index + 1 < exerciseState.sentence.size) {
                                    append(" ")
                                    withStyle(style = SpanStyle(textDecoration = TextDecoration.Underline)) {
                                        append(missedWord)
                                    }
                                    append(" ")
                                }
                            }
                        }
                    },
                    fontSize = 20.sp
                )
            }

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                stringResource(R.string.choose_missing_word),
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
                        onClick = { fillGapsObserver.onVariantClick(index) },
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

        Box(
            modifier = Modifier.fillMaxWidth().height(160.dp),
            contentAlignment = Alignment.Center
        ) {
            if (isCompleted) {
                ExerciseContinueButton(
                    onClick = { fillGapsObserver.onCompleteExercise() }
                )
            }
        }
    }
}

@Preview(device = PIXEL_7)
@DevicePreviews
@Composable
fun FillGapsExercisePreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface)
        ) {
            FillGapsExercise(
                modifier = Modifier.fillMaxSize(),
                exerciseState = Exercise.FillTheGap(
                    sentence = listOf("I build a", "last year"),
                    answerVariants = listOf("Car", "Bird", "Brother", "House"),
                    correctVariant = 3,
                    selectedVariant = 2
                ),
                fillGapsObserver = object : FillGapsObserver() {
                    override fun onCompleteExercise() {}
                    override fun onVariantClick(variantIndex: Int) {}
                },
                isCompleted = true
            )
        }
    }
}


@Preview(device = PIXEL_7)
@Composable
fun FillGapsExerciseLongPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface)
        ) {
            FillGapsExercise(
                modifier = Modifier.fillMaxSize(),
                exerciseState = Exercise.FillTheGap(
                    sentence = listOf(
                        "I build a",
                        "last year and very long and Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
                    ),
                    answerVariants = listOf("Car", "Bird", "Brother", "House"),
                    correctVariant = 3,
                    selectedVariant = null
                ),
                fillGapsObserver = object : FillGapsObserver() {
                    override fun onCompleteExercise() {}
                    override fun onVariantClick(variantIndex: Int) {}
                },
                isCompleted = true
            )
        }
    }
}
