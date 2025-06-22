package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.ui.components.TopAppBar
import ru.fluentlyapp.fluently.ui.screens.lesson.components.LoadingLessonComponent
import ru.fluentlyapp.fluently.ui.screens.lesson.components.exercises.Exercise
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun LessonScreen(
    modifier: Modifier = Modifier,
    lessonScreenViewModel: LessonScreenViewModel = hiltViewModel(),
    onBackClick: () -> Unit
) {
    val uiState by lessonScreenViewModel.uiState.collectAsState()

    Column(
        modifier = modifier.background(color = FluentlyTheme.colors.surface),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        TopAppBar(
            modifier = Modifier.fillMaxWidth(),
            onBackClick = onBackClick
        )
        LessonScreenContent(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth(),
            uiState = uiState,
            onContinue = { lessonScreenViewModel.moveToNextComponent() },
            onUpdateComponent = { lessonScreenViewModel.updateCurrentComponent(it) }
        )
    }
}

@Composable
fun LessonScreenContent(
    modifier: Modifier = Modifier,
    uiState: LessonScreenUiState,
    onContinue: () -> Unit,
    onUpdateComponent: (newComponent: LessonComponent) -> Unit,
) {
    Column(
        modifier = modifier.background(color = FluentlyTheme.colors.surface).padding(horizontal = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {

        Box(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
        ) {
            if (uiState.currentComponent is Exercise) {
                Exercise(
                    modifier = Modifier.fillMaxSize(),
                    exercise = uiState.currentComponent,
                    onUpdateExercise = onUpdateComponent
                )
            } else if (uiState.currentComponent is LessonComponent.Loading) {
                LoadingLessonComponent(modifier = Modifier.fillMaxSize())
            }
        }

        if (uiState.showContinueButton) {
            Box(
                modifier = Modifier
                    .clickable(onClick = onContinue)
                    .padding(bottom = 32.dp)
                    .clip(RoundedCornerShape(16.dp))
                    .background(color = FluentlyTheme.colors.primary)
                    .fillMaxWidth(.7f)
                    .height(80.dp)
            ) {
                Text(
                    modifier = Modifier.align(Alignment.Center),
                    text = stringResource(R.string.continue_to_next_exercise),
                    fontSize = 24.sp,
                    fontWeight = FontWeight.Bold,
                    color = FluentlyTheme.colors.onPrimary
                )
            }
        } else {
            Box(
                modifier = Modifier
                    .padding(bottom = 32.dp)
                    .height(80.dp)
                    .fillMaxWidth()
            )
        }
    }
}

@Preview(device = PIXEL_7)
@Composable
fun LessonScreenPreview() {
    FluentlyTheme {
        LessonScreenContent(
            modifier = Modifier.background(color = FluentlyTheme.colors.surface).padding(horizontal = 16.dp),
            uiState = LessonScreenUiState(
                currentComponent = Exercise.ChooseTranslation(
                    word = "Influence",
                    answerVariants = listOf(
                        "Влияние", "Благодарность", "Двойственность", "Комар"
                    ),
                    correctVariant = 0,
                    selectedVariant = 2,
                ),
                showContinueButton = true
            ),
            onContinue = {},
            onUpdateComponent = {}
        )
    }
}

@Preview(device = PIXEL_7)
@Composable
fun LessonScreen2Preview() {
    FluentlyTheme {
        LessonScreenContent(
            uiState = LessonScreenUiState(
                currentComponent = Exercise.NewWord(
                    word = "awareness",
                    translation = "осознание",
                    phoneticTranscription = "/əˈweə.nəs/",
                    doesUserKnow = null,
                    examples = listOf("Environmental awareness is rising" to "Экологическое осознание растет")
                ),
                showContinueButton = true
            ),
            onContinue = {},
            onUpdateComponent = {}
        )
    }
}
