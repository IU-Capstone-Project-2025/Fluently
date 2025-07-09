package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.ui.components.TopAppBar

@Composable
fun LessonFlowScreen(
    modifier: Modifier = Modifier,
    lessonFlowViewModel: LessonFlowViewModel = hiltViewModel(),
    onBackClick: () -> Unit
) {
    val currentComponent by lessonFlowViewModel.currentComponent.collectAsState()

    LaunchedEffect(onBackClick) {
        withContext(Dispatchers.Main.immediate) {
            lessonFlowViewModel.commandsChannel.collect { command ->
                if (command is LessonFlowCommand.UserFinishesLesson) {
                    onBackClick()
                }
            }
        }
    }

    Column(modifier = modifier) {
        TopAppBar(modifier = Modifier.fillMaxWidth(), onBackClick = onBackClick)

        LessonComponentRenderer(
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f),
            component = LessonComponentWithIndex(
                currentComponent ?: Decoration.Loading(),
                currentComponent?.id ?: -1
            ),
            chooseTranslationObserver = lessonFlowViewModel.chooseTranslationObserver,
            newWordObserver = lessonFlowViewModel.newWordObserver,
            fillGapsObserver = lessonFlowViewModel.fillGapsObserver,
            inputWordObserver = lessonFlowViewModel.inputWordObserver,
            onboardingDecorationObserver = lessonFlowViewModel.onboardingDecorationObserver,
            finishDecorationObserver = lessonFlowViewModel.finishDecorationObserver
        )
    }
}
