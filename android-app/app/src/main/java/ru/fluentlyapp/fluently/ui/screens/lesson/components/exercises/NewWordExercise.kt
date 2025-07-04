package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.ui.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

abstract class NewWordObserver {
    abstract fun onUserKnowsWord(doesUserKnowWord: Boolean)
    abstract fun onCompleteExercise()
}

@Composable
fun NewWordExercise(
    modifier: Modifier = Modifier,
    exerciseState: Exercise.NewWord,
    newWordObserver: NewWordObserver,
    isCompleted: Boolean
) {
    Column(
        modifier = modifier.background(FluentlyTheme.colors.surface),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Column(
            modifier = Modifier.weight(1f).fillMaxWidth(.8f),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Column(
                modifier = Modifier
                    .verticalScroll(state = rememberScrollState())
                    .fillMaxWidth()
                    .heightIn(max = 600.dp)
                    .clip(RoundedCornerShape(size = 16.dp))
                    .background(color = FluentlyTheme.colors.surfaceContainerHigh)
                    .padding(16.dp)
            ) {
                Text(
                    text = exerciseState.word,
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onSurface
                )
                Spacer(modifier = Modifier.height(16.dp))
                Text(
                    text = stringResource(R.string.translation),
                    color = FluentlyTheme.colors.onSurfaceVariant
                )
                Text(
                    text = exerciseState.translation
                )
                Spacer(modifier = Modifier.height(16.dp))

                Text(
                    text = "Примеры",
                    color = FluentlyTheme.colors.onSurfaceVariant
                )
                repeat(exerciseState.examples.size) { index ->
                    val (english, translation) = exerciseState.examples[index]
                    Text(text = english)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = translation)
                    if (index != exerciseState.examples.size - 1) {
                        Box(
                            modifier = Modifier
                                .height(16.dp)
                                .fillMaxWidth(),
                            contentAlignment = Alignment.Center
                        ) {
                            Spacer(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(1.dp)
                                    .border(width = 1.dp, color = FluentlyTheme.colors.onSurfaceVariant)
                            )
                        }
                    }
                }
            }

            Spacer(modifier = Modifier.height(16.dp))

            Row(modifier = Modifier.fillMaxWidth()) {
                Box(
                    modifier = Modifier
                        .alpha(if (exerciseState.doesUserKnow == false) .3f else 1f)
                        .clip(RoundedCornerShape(12.dp))
                        .clickable(
                            enabled = !isCompleted,
                            onClick = { newWordObserver.onUserKnowsWord(true) }
                        )
                        .border(
                            color = FluentlyTheme.colors.onSurface,
                            width = 2.dp,
                            shape = RoundedCornerShape(12.dp)
                        )
                        .weight(1f)
                        .padding(16.dp)
                ) {
                    Text(
                        modifier = Modifier.fillMaxWidth(),
                        text = "ЗНАЮ",
                        fontWeight = FontWeight.Bold,
                        textAlign = TextAlign.Center,
                        color = FluentlyTheme.colors.onSurfaceVariant
                    )
                }

                Spacer(modifier = Modifier.width(16.dp))

                Box(
                    modifier = Modifier
                        .alpha(if (exerciseState.doesUserKnow == true) .3f else 1f)
                        .clip(RoundedCornerShape(12.dp))
                        .clickable(
                            enabled = !isCompleted,
                            onClick = { newWordObserver.onUserKnowsWord(false) }
                        )
                        .weight(1f)
                        .background(color = FluentlyTheme.colors.secondary)
                        .padding(16.dp),
                ) {
                    Text(
                        modifier = Modifier.fillMaxWidth(),
                        text = "УЧИТЬ",
                        fontWeight = FontWeight.Bold,
                        textAlign = TextAlign.Center,
                        color = FluentlyTheme.colors.onPrimary
                    )
                }
            }
        }

        ExerciseContinueButton(
            modifier = Modifier.fillMaxWidth(.7f).padding(32.dp).height(80.dp),
            enabled = isCompleted,
            onClick = newWordObserver::onCompleteExercise
        )
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun NewWordExercisePreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(color = FluentlyTheme.colors.surface),
            contentAlignment = Alignment.Center
        ) {
            NewWordExercise(
                modifier = Modifier.fillMaxSize(),
                exerciseState = Exercise.NewWord(
                    word = "Deprecation",
                    phoneticTranscription = "/ˌdep.rəˈkeɪ.ʃən/",
                    doesUserKnow = true,
                    translation = "Устеревание",
                    examples = listOf(
                        "This function is deprecated since lirbary version 1.2" to
                                "Эта функция устарела, начиная с  версии билиотеки 1.2",
                        "Components Deprecation is a main source of conflicts in androidandroid" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                    )
                ),
                newWordObserver = object : NewWordObserver() {
                    override fun onUserKnowsWord(doesUserKnowWord: Boolean) {}
                    override fun onCompleteExercise() {}
                },
                isCompleted = true
            )
        }
    }
}
