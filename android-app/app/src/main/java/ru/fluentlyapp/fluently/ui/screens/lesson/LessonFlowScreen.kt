package ru.fluentlyapp.fluently.ui.screens.lesson

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.ui.components.TopAppBar

@Composable
fun LessonFlowScreen(
    modifier: Modifier = Modifier,
    lessonFlowViewModel: LessonFlowViewModel = hiltViewModel(),
    onBackClick: () -> Unit
) {
    val lesson by lessonFlowViewModel.lesson.collectAsState()
    val currentComponent = lesson?.currentComponent ?: LessonComponent.Loading

    Column(modifier = modifier) {
        TopAppBar(modifier = Modifier.fillMaxWidth(), onBackClick = onBackClick)

        LessonComponentRenderer(
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f),
            component = LessonComponentWithIndex(
                currentComponent,
                lesson?.currentLessonComponentIndex ?: -1
            ),
            chooseTranslationObserver = lessonFlowViewModel.chooseTranslationObserver,
            newWordObserver = lessonFlowViewModel.newWordObserver,
            fillGapsObserver = lessonFlowViewModel.fillGapsObserver,
            inputWordObserver = lessonFlowViewModel.inputWordObserver
        )
    }
}
