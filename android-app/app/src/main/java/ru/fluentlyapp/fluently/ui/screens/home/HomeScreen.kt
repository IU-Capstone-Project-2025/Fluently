package ru.fluentlyapp.fluently.ui.screens.home

import androidx.compose.animation.animateContentSize
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.heightIn
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.systemBars
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.widthIn
import androidx.compose.foundation.layout.windowInsetsPadding
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.Hyphens
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil3.compose.AsyncImage
import coil3.request.ImageRequest
import coil3.request.error
import coil3.request.placeholder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.previewdata.words
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreenUiState.OngoingLessonState
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

@Composable
fun HomeScreen(
    modifier: Modifier = Modifier,
    homeScreenViewModel: HomeScreenViewModel = hiltViewModel(),
    onNavigateToLesson: () -> Unit,
    onNavigateToCalendar: () -> Unit,
    onLearnedWordsClick: () -> Unit,
    onInProgressWordsClick: () -> Unit,
) {
    val uiState by homeScreenViewModel.uiState.collectAsState()

    LaunchedEffect(onNavigateToLesson) {
        withContext(Dispatchers.Main.immediate) {
            homeScreenViewModel.commandsChannel.collect { command ->
                if (command is HomeCommands.NavigateToLesson) {
                    onNavigateToLesson()
                }
            }
        }
    }

    HomeScreenContent(
        modifier = modifier,
        uiState = uiState,
        onLessonClick = {
            homeScreenViewModel.ensureOngoingLesson()
        },
        onCalendarClick = onNavigateToCalendar,
        onLearnedWordsClick = onLearnedWordsClick,
        onInProgressWordsClick = onInProgressWordsClick,
        onStartLearningWordOfTheDay = {
            homeScreenViewModel.startLearningWordOfTheDay()
        }
    )
}

@Composable
fun HomeScreenContent(
    modifier: Modifier = Modifier,
    uiState: HomeScreenUiState,
    onLessonClick: () -> Unit,
    onCalendarClick: () -> Unit,
    onLearnedWordsClick: () -> Unit,
    onInProgressWordsClick: () -> Unit,
    onStartLearningWordOfTheDay: () -> Unit,
) {
    Column(
        modifier = modifier
            .background(color = FluentlyTheme.colors.primary)
            .windowInsetsPadding(WindowInsets.systemBars)
    ) {
        Row(
            modifier = Modifier
                .height(160.dp)
                .fillMaxWidth()
                .padding(horizontal = 32.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = "Домашний экран",
                    lineHeight = 40.sp,
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onPrimary,
                    fontWeight = FontWeight.Bold
                )
            }
            AsyncImage(
                model = ImageRequest.Builder(LocalContext.current)
                    .data(uiState.avatarPicture)
                    .placeholder(R.drawable.ic_funny_square)
                    .error(R.drawable.ic_funny_square)
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
                    .background(FluentlyTheme.colors.surfaceContainerHigh)
            )
        }
        Column(
            modifier = Modifier
                .clip(RoundedCornerShape(topStart = 32.dp, topEnd = 32.dp))
                .fillMaxWidth()
                .weight(1f)
                .background(color = FluentlyTheme.colors.surface)
                .padding(vertical = 32.dp, horizontal = 20.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.SpaceEvenly
        ) {
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Column(
                    modifier = Modifier
                        .clip(RoundedCornerShape(16.dp))
                        .background(FluentlyTheme.colors.surfaceInverse)
                        .padding(16.dp)
                        .animateContentSize(),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    if (uiState.wordOfTheDay == null) {
                        CircularProgressIndicator(color = FluentlyTheme.colors.onSurfaceInverse)
                    } else {
                        Text(
                            fontSize = 32.sp,
                            color = FluentlyTheme.colors.onSurfaceInverse,
                            text = uiState.wordOfTheDay.word
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            fontSize = 20.sp,
                            color = FluentlyTheme.colors.onSurfaceVariantInverse,
                            text = uiState.wordOfTheDay.translation
                        )
                    }
                }
                Spacer(modifier = Modifier.height(4.dp))
                Row(
                    modifier = Modifier
                        .alpha(if (uiState.wordOfTheDay == null) .5f else 1f)
                        .clip(RoundedCornerShape(16.dp))
                        .clickable(
                            enabled = uiState.wordOfTheDay != null && !uiState.hasWordOfTheDaySaved,
                            onClick = onStartLearningWordOfTheDay
                        )
                        .background(FluentlyTheme.colors.surfaceContainerHigh)
                        .padding(horizontal = 16.dp, vertical = 8.dp)
                        .animateContentSize()
                ) {
                    if (uiState.wordOfTheDay == null || !uiState.hasWordOfTheDaySaved) {
                        Icon(
                            modifier = Modifier.size(24.dp),
                            painter = painterResource(R.drawable.ic_plus_circle),
                            tint = FluentlyTheme.colors.onSurface,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text("Изучать")
                    } else {
                        Icon(
                            modifier = Modifier.size(24.dp),
                            painter = painterResource(R.drawable.ic_check_circle),
                            tint = FluentlyTheme.colors.onSurface,
                            contentDescription = null
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                        Text("Сохранено!")
                    }
                }
            }

            Row(modifier = Modifier.fillMaxWidth()) {
                Column(
                    modifier = Modifier
                        .heightIn(min = 130.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .background(color = FluentlyTheme.colors.primaryVariant)
                        .clickable(onClick = onCalendarClick)
                        .padding(16.dp)
                        .weight(1f),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.SpaceEvenly
                ) {
                    Icon(
                        modifier = Modifier.size(60.dp),
                        painter = painterResource(R.drawable.ic_calendar),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.primary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = "Календарь",
                        fontSize = 12.sp,
                        softWrap = true,
                        style = TextStyle(
                            hyphens = Hyphens.Unspecified
                        )
                    )
                }
                Spacer(modifier = Modifier.width(12.dp))
                Column(
                    modifier = Modifier
                        .heightIn(min = 130.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .clickable(onClick = onLearnedWordsClick)
                        .background(color = FluentlyTheme.colors.secondaryVariant)
                        .padding(16.dp)
                        .weight(1f)
                ) {
                    Icon(
                        modifier = Modifier.size(24.dp),
                        painter = painterResource(R.drawable.ic_person_learned_words),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.secondary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = uiState.learnedWordsNumber.toString(), fontSize = 32.sp)
                    Text(text = "Изучено", fontSize = 12.sp)
                }
                Spacer(modifier = Modifier.width(12.dp))
                Column(
                    modifier = Modifier
                        .heightIn(min = 120.dp)
                        .clip(RoundedCornerShape(16.dp))
                        .clickable(onClick = onInProgressWordsClick)
                        .background(color = FluentlyTheme.colors.tertiaryVariant1)
                        .padding(16.dp)
                        .weight(1f)
                ) {
                    Icon(
                        modifier = Modifier.size(24.dp),
                        painter = painterResource(R.drawable.ic_progress),
                        contentDescription = null,
                        tint = FluentlyTheme.colors.tertiary
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = uiState.inProgressWordsNumber.toString(), fontSize = 32.sp)
                    Text(text = "Изучаются", fontSize = 12.sp)
                }
            }
            Box(
                modifier = Modifier
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

@DevicePreviews
@Composable
fun HomeScreenPreview() {
    FluentlyTheme {
        HomeScreenContent(
            modifier = Modifier.fillMaxSize(),
            uiState = HomeScreenUiState(
                ongoingLessonState = OngoingLessonState.LOADING,
                wordOfTheDay = words[0]
            ),
            onLessonClick = {},
            onCalendarClick = { },
            onLearnedWordsClick = {},
            onInProgressWordsClick = {},
            onStartLearningWordOfTheDay = {}
        )
    }
}