package ru.fluentlyapp.fluently.ui.screens.calendar

import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.animateColorAsState
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews
import androidx.compose.runtime.getValue
import ru.fluentlyapp.fluently.ui.previewdata.words
import ru.fluentlyapp.fluently.ui.theme.components.WordList
import java.time.LocalDate

@Composable
fun CalendarScreen(
    modifier: Modifier = Modifier,
    onBackClick: () -> Unit,
    calendarViewModel: CalendarViewModel = hiltViewModel()
) {
    val uiState by calendarViewModel.uiState.collectAsState()
    CalendarScreenContent(
        modifier = modifier,
        uiState = uiState,
        onBackClick = onBackClick,
        onDayClick = { calendarViewModel.setDate(it) },
        onShowInProgressWords = { showInProgress ->
            calendarViewModel.updateShowIsLearning(showInProgress)
        }
    )
}

@Composable
fun CalendarScreenContent(
    modifier: Modifier = Modifier,
    uiState: CalendarScreenUiState,
    onDayClick: (LocalDate) -> Unit,
    onBackClick: () -> Unit,
    onShowInProgressWords: (Boolean) -> Unit
) {
    Column(
        modifier = modifier.background(color = FluentlyTheme.colors.primary)
    ) {
        Box(
            modifier = Modifier
                .height(80.dp)
                .fillMaxWidth()
                .padding(horizontal = 8.dp),
        ) {
            Icon(
                modifier = Modifier
                    .clip(CircleShape)
                    .clickable(onClick = onBackClick)
                    .align(Alignment.CenterStart),
                tint = FluentlyTheme.colors.primaryVariant,
                painter = painterResource(R.drawable.ic_chevron_left),
                contentDescription = "Back button"
            )

            Text(
                modifier = Modifier.align(Alignment.Center),
                text = "Календарь",
                fontSize = 28.sp,
                color = FluentlyTheme.colors.onPrimary,
                fontWeight = FontWeight.Bold,
            )
        }
        Column(
            modifier = Modifier
                .clip(RoundedCornerShape(topStart = 32.dp, topEnd = 32.dp))
                .fillMaxWidth()
                .weight(1f)
                .background(color = FluentlyTheme.colors.surface)
                .padding(top = 20.dp, start = 20.dp, end = 20.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Top
        ) {
            CalendarScroll(
                modifier = Modifier.fillMaxWidth(),
                selectedDay = uiState.day,
                onDayClick = onDayClick
            )

            Row(
                modifier = Modifier
                    .padding(vertical = 8.dp)
                    .fillMaxWidth()
                    .height(40.dp),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                SwitchText(
                    text = "Выученные слова",
                    isActive = !uiState.showIsLearning,
                    onClick = { onShowInProgressWords(false) }
                )
                SwitchText(
                    text = "Слова в обучении",
                    isActive = uiState.showIsLearning,
                    onClick = { onShowInProgressWords(true) }
                )
            }

            AnimatedContent(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxWidth(),
                targetState = uiState.showIsLearning
            ) { targetIsLearning ->
                if (targetIsLearning) {
                    WordList(
                        modifier = Modifier.fillMaxSize(),
                        words = uiState.inProgressWords
                    )
                } else {
                    WordList(
                        modifier = Modifier.fillMaxSize(),
                        words = uiState.learnedWords
                    )
                }
            }
        }
    }
}

@Composable
fun SwitchText(
    modifier: Modifier = Modifier,
    text: String,
    isActive: Boolean,
    onClick: () -> Unit
) {
    val textColor by animateColorAsState(
        targetValue = if (isActive)
            FluentlyTheme.colors.onSurface
        else
            FluentlyTheme.colors.onSurfaceVariant
    )
    Box(
        modifier = modifier
            .fillMaxHeight()
            .clip(RoundedCornerShape(100.dp))
            .clickable(onClick = onClick)
            .padding(horizontal = 8.dp),
        contentAlignment = Alignment.Center
    ) {
        Text(
            text = text,
            color = textColor
        )
    }
}


@DevicePreviews
@Composable
fun CalendarScreenPreview() {
    FluentlyTheme {
        CalendarScreenContent(
            uiState = CalendarScreenUiState(
                day = LocalDate.now(),
                learnedWords = (1..50).map { words[0] },
                inProgressWords = (1..50).map { words[1] },
                showIsLearning = true
            ),
            onBackClick = {},
            onDayClick = {},
            onShowInProgressWords = {}
        )
    }
}
