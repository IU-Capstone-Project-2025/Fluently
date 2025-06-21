package ru.fluentlyapp.fluently.ui.screens.home

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import androidx.compose.runtime.getValue

@Composable
fun HomeScreen(
    modifier: Modifier = Modifier,
    homeScreenViewModel: HomeScreenViewModel = hiltViewModel(),
    onNavigateToLesson: () -> Unit
) {
    val uiState by homeScreenViewModel.uiState.collectAsState()
    HomeScreenContent(
        modifier = modifier,
        uiState = uiState,
        onLessonClick = onNavigateToLesson
    )
}

@Composable
fun HomeScreenContent(
    modifier: Modifier = Modifier,
    uiState: HomeScreenUiState,
    onLessonClick: () -> Unit
) {
    Box(modifier = modifier.background(color = FluentlyTheme.colors.surface)) {
        Button(
            modifier = Modifier.align(Alignment.Center),
            onClick = onLessonClick,
            colors = ButtonDefaults.buttonColors(containerColor = FluentlyTheme.colors.primary)
        ) {
            Text(
                text = if (uiState.hasOngoingLesson) "Continue Lesson" else "Start New Lesson"
            )
        }
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun HomeScreenPreview() {
    FluentlyTheme {
        HomeScreenContent(
            modifier = Modifier.fillMaxSize(),
            uiState = HomeScreenUiState(hasOngoingLesson = false),
            onLessonClick = {}
        )
    }
}