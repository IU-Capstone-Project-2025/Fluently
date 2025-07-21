package ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.ui.theme.components.ExerciseContinueButton
import ru.fluentlyapp.fluently.ui.theme.components.NewWordCard
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.MediumPhonePreview
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview

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
        modifier = modifier
            .background(FluentlyTheme.colors.surface)
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Box(modifier = Modifier.weight(1f), contentAlignment = Alignment.Center) {
            NewWordCard(
                modifier = Modifier.fillMaxWidth(.8f),
                word = exerciseState.word,
                translation = exerciseState.translation,
                examples = exerciseState.examples
            )
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
                    text = stringResource(R.string.know),
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
                    text = stringResource(R.string.study_capitalized),
                    fontWeight = FontWeight.Bold,
                    textAlign = TextAlign.Center,
                    color = FluentlyTheme.colors.onPrimary
                )
            }
        }
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(120.dp),
            contentAlignment = Alignment.Center
        ) {
            if (isCompleted) {
                ExerciseContinueButton(
                    onClick = { newWordObserver.onCompleteExercise() }
                )
            }
        }
    }

}

@SmallPhonePreview
@MediumPhonePreview
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
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                    ),
                    wordId = ""
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

@SmallPhonePreview
@MediumPhonePreview
@Composable
fun NewWordExerciseScrollPreview() {
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
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                        "Components Deprecation is a main source of conflicts in android" to
                                "Устаревание компонентов - главная причина конфликтов в Андроиде",
                    ),
                    wordId = ""
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
