package ru.fluentlyapp.fluently.ui.screens.home

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
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil3.compose.AsyncImage
import coil3.request.ImageRequest
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreenUiState.OngoingLessonState

@Composable
fun HomeScreen(
    modifier: Modifier = Modifier,
    homeScreenViewModel: HomeScreenViewModel = hiltViewModel(),
    onNavigateToLesson: (lessonId: String) -> Unit
) {
    val uiState by homeScreenViewModel.uiState.collectAsState()
    val ongoingLessonIsReady by homeScreenViewModel.ongoingLessonIsReady.collectAsState()

    LaunchedEffect(ongoingLessonIsReady) {
        if (ongoingLessonIsReady) {
            onNavigateToLesson("")
        }
    }

    HomeScreenContent(
        modifier = modifier,
        uiState = uiState,
        onLessonClick = {
            homeScreenViewModel.ensureOngoingLesson()
        }
    )
}

@Composable
fun HomeScreenContent(
    modifier: Modifier = Modifier,
    uiState: HomeScreenUiState,
    onLessonClick: () -> Unit
) {
    Column(
        modifier = modifier.background(color = FluentlyTheme.colors.primary)
    ) {
        Row(
            modifier = Modifier
                .height(200.dp)
                .fillMaxWidth()
                .padding(horizontal = 32.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = "Цель:",
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onPrimary,
                    fontWeight = FontWeight.Bold
                )
                Text(
                    text = uiState.goal,
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onPrimary,
                    fontWeight = FontWeight.Bold
                )
            }
            AsyncImage(
                model = ImageRequest.Builder(LocalContext.current)
                    .data(uiState.avatarPicture)
                    .build(),
                contentScale = ContentScale.Crop,
                contentDescription = "Avatar Picture",
                modifier = Modifier
                    .clip(CircleShape)
                    .border(
                        width = 2.dp,
                        shape = CircleShape,
                        color = FluentlyTheme.colors.onPrimary
                    )
                    .size(88.dp)
                    .background(Color.Black)
            )
        }
        Column(
            modifier = Modifier
                .clip(RoundedCornerShape(topStart = 32.dp, topEnd = 32.dp))
                .fillMaxWidth()
                .weight(1f)
                .background(color = FluentlyTheme.colors.surface)
                .padding(vertical = 32.dp, horizontal = 20.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                text = "Word of the day",
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )

            Spacer(modifier = Modifier.height(8.dp))

            Column(
                modifier = Modifier
                    .clip(shape = RoundedCornerShape(16.dp))
                    .background(color = FluentlyTheme.colors.surfaceInverse)
                    .padding(vertical = 20.dp, horizontal = 32.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = uiState.wordOfTheDay,
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onSurfaceInverse
                )
                Text(
                    text = uiState.wordOfTheDayTranslation,
                    color = FluentlyTheme.colors.onSurfaceVariantInverse
                )
            }

            Spacer(modifier = Modifier.height(4.dp))

            Row(
                modifier = Modifier
                    .clip(RoundedCornerShape(100.dp))
                    .background(color = FluentlyTheme.colors.surfaceContainerHigh)
                    .padding(4.dp),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                Icon(
                    modifier = Modifier.size(20.dp),
                    painter = painterResource(R.drawable.ic_plus_circle),
                    contentDescription = null,
                    tint = FluentlyTheme.colors.onSurface
                )
                Spacer(modifier = Modifier.width(4.dp))
                Text(text = "Добавить в коллекцию")
            }

            Spacer(modifier = Modifier.height(16.dp))

            Row(modifier = Modifier.fillMaxWidth()) {
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(16.dp))
                        .background(color = FluentlyTheme.colors.primaryVariant)
                        .padding(16.dp)
                        .weight(1f)
                ) {
                    Icon(
                        modifier = Modifier.size(24.dp),
                        painter = painterResource(R.drawable.ic_book),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.primary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = uiState.notesNumber.toString(), fontSize = 32.sp)
                    Text(text = "Заметок", fontSize = 14.sp)
                }
                Spacer(modifier = Modifier.width(12.dp))
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(16.dp))
                        .background(color = FluentlyTheme.colors.secondaryVariant)
                        .padding(16.dp)
                        .weight(1f)
                ) {
                    Icon(
                        modifier = Modifier.size(24.dp),
                        painter = painterResource(R.drawable.ic_bachelor_hat),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.secondary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = uiState.learnedWordsNumber.toString(), fontSize = 32.sp)
                    Text(text = "Изучено", fontSize = 14.sp)
                }
                Spacer(modifier = Modifier.width(12.dp))
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(16.dp))
                        .background(color = FluentlyTheme.colors.tertiaryVariant1)
                        .padding(16.dp)
                        .weight(1f)
                ) {
                    Icon(
                        modifier = Modifier.size(24.dp),
                        painter = painterResource(R.drawable.ic_bachelor_hat),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.tertiary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = uiState.notesNumber.toString(), fontSize = 32.sp)
                    Text(text = "Не изучено", fontSize = 14.sp)
                }
            }
            Box(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth()
            ) {
                Row(
                    modifier = Modifier
                        .clip(RoundedCornerShape(100.dp))
                        .clickable(
                            onClick = onLessonClick,
                            enabled = uiState.ongoingLessonState != OngoingLessonState.LOADING
                        )
                        .alpha(if (uiState.ongoingLessonState == OngoingLessonState.LOADING) .5f else 1f)
                        .background(color = FluentlyTheme.colors.surfaceInverse)
                        .padding(12.dp)
                        .height(40.dp)
                        .widthIn(min = 240.dp)
                        .align(Alignment.Center),
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.SpaceEvenly
                ) {
                    if (uiState.ongoingLessonState == OngoingLessonState.LOADING) {
                        CircularProgressIndicator(color = FluentlyTheme.colors.onSurfaceInverse)
                        Spacer(modifier = Modifier.width(16.dp))
                    }
                    Text(
                        text = when (uiState.ongoingLessonState) {
                            OngoingLessonState.ERROR -> "Ошибка :("
                            OngoingLessonState.HAS_PAUSED -> "Продолжить урок"
                            OngoingLessonState.LOADING -> "Загружаем урок..."
                            OngoingLessonState.NOT_STARTED -> "Начать урок"
                        },
                        fontSize = 24.sp,
                        color = FluentlyTheme.colors.onSurfaceInverse
                    )
                }
            }
        }
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun HomeScreenPreview() {
    FluentlyTheme {
        HomeScreenContent(
            modifier = Modifier.fillMaxSize(),
            uiState = HomeScreenUiState(ongoingLessonState = OngoingLessonState.LOADING),
            onLessonClick = { }
        )
    }
}