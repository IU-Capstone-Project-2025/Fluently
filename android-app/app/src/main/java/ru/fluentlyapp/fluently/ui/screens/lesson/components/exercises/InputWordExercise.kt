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
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.saveable.rememberSaveable
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.SpanStyle
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.buildAnnotatedString
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.withStyle
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.ui.theme.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

interface InputWordObserver {
    fun onCompleteExercise()
    fun onConfirmInput(input: String)
}

@Composable
fun InputWordExercise(
    modifier: Modifier,
    exerciseState: Exercise.InputWord,
    observer: InputWordObserver,
    isCompleted: Boolean
) {
    Column(
        modifier = modifier.background(FluentlyTheme.colors.surface)
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f)
                .verticalScroll(state = rememberScrollState())
                .padding(horizontal = 16.dp),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                text = exerciseState.translation,
                fontSize = 32.sp
            )

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                text = stringResource(R.string.enter_the_word_translation),
                fontSize = 16.sp
            )
            Spacer(modifier = Modifier.height(8.dp))
            var inputWordValue by rememberSaveable { mutableStateOf(exerciseState.inputtedWord ?: "") }

            val inputFieldModifier = Modifier.border(
                shape = RoundedCornerShape(8.dp),
                width = 2.dp,
                color = if (exerciseState.isAnswered) {
                    if (exerciseState.correctAnswer == exerciseState.inputtedWord) {
                        FluentlyTheme.colors.correct
                    } else {
                        FluentlyTheme.colors.wrong
                    }
                } else {
                    Color.Unspecified
                }
            )
            BasicTextField(
                value = inputWordValue,
                enabled = !isCompleted,
                onValueChange = { inputWordValue = it },
                textStyle = TextStyle(fontSize = 16.sp),
                singleLine = true,
                decorationBox = { innerTextField ->
                    Box(
                        modifier = Modifier
                            .clip(shape = RoundedCornerShape(8.dp))
                            .background(FluentlyTheme.colors.surface)
                            then inputFieldModifier
                            .padding(8.dp)
                            .fillMaxWidth(.8f)
                    ) {
                        innerTextField()
                    }
                }
            )


            if (exerciseState.isAnswered && exerciseState.inputtedWord != exerciseState.correctAnswer) {
                Text(
                    text = buildAnnotatedString {
                        append("Правильный ответ: ")
                        withStyle(SpanStyle(fontWeight = FontWeight.Bold)) {
                            append(exerciseState.correctAnswer)
                        }
                    }
                )
            }

            Spacer(modifier = Modifier.height(16.dp))

            Box(
                modifier = Modifier
                    .alpha(if (inputWordValue.isEmpty() || isCompleted) .5f else 1f)
                    .clip(RoundedCornerShape(8.dp))
                    .background(FluentlyTheme.colors.secondary)
                    .clickable(
                        enabled = inputWordValue.isNotEmpty() && !isCompleted,
                        onClick = { observer.onConfirmInput(inputWordValue) }
                    )
                    .padding(16.dp)
            ) {
                Text(
                    text = stringResource(R.string.check),
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Bold,
                    color = FluentlyTheme.colors.onSecondary
                )
            }
        }
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(160.dp),
            contentAlignment = Alignment.Center
        ) {
            if (isCompleted) {
                ExerciseContinueButton(
                    onClick = { observer.onCompleteExercise() }
                )
            }
        }
    }
}

@DevicePreviews
@Composable
fun InputWordExercisePreview() {
    FluentlyTheme {
        InputWordExercise(
            modifier = Modifier.fillMaxSize(),
            exerciseState = Exercise.InputWord("Устаревание", wordId = "", "Deprecation", "Deprecation"),
            observer = object : InputWordObserver {
                override fun onCompleteExercise() {}
                override fun onConfirmInput(input: String) {}
            },
            isCompleted = true
        )
    }
}