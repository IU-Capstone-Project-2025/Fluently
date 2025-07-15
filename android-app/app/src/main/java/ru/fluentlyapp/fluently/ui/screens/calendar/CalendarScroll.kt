package ru.fluentlyapp.fluently.ui.screens.calendar

import androidx.compose.animation.animateColorAsState
import androidx.compose.animation.rememberSplineBasedDecay
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.gestures.FlingBehavior
import androidx.compose.foundation.gestures.ScrollScope
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.derivedStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews
import java.time.LocalDate
import java.time.YearMonth
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import java.util.Locale


@Composable
fun CalendarScroll(
    modifier: Modifier,
    selectedDay: LocalDate,
    onDayClick: (LocalDate) -> Unit,
) {
    val formatter = remember { DateTimeFormatter.ofPattern("MMM yyyy", Locale.ENGLISH) }
    val beginDate = remember { LocalDate.of(1970, 1, 1) }
    val currentDate = remember { LocalDate.now() }
    val numberOfDays = remember { ChronoUnit.DAYS.between(beginDate, currentDate).toInt() + 1 }

    val listState = rememberLazyListState()

    LaunchedEffect(Unit) {
        listState.scrollToItem(numberOfDays)
    }

    val yearMonth by remember {
        derivedStateOf {
            YearMonth.from(beginDate.plusDays(listState.firstVisibleItemIndex.toLong()))
        }
    }

    Column(modifier = modifier.fillMaxWidth(), horizontalAlignment = Alignment.CenterHorizontally) {
        Text(
            text = yearMonth.format(formatter),
            fontWeight = FontWeight.Bold,
            color = FluentlyTheme.colors.secondary
        )

        LazyRow(
            modifier = Modifier.fillMaxWidth(),
            state = listState,
        ) {
            items(numberOfDays) { index ->
                val date = beginDate.plusDays(index.toLong())
                CalendarDay(
                    modifier = Modifier.size(50.dp),
                    date = date,
                    isSelected = date == selectedDay,
                    onClick = { onDayClick(date) },
                    startSeparator = date.dayOfMonth == 1
                )
            }
        }

    }
}

@Composable
fun CalendarDay(
    modifier: Modifier = Modifier,
    date: LocalDate,
    isSelected: Boolean,
    onClick: () -> Unit,
    endSeparator: Boolean = false,
    startSeparator: Boolean = false,
) {
    val selectionColor by animateColorAsState(
        targetValue = if (isSelected) FluentlyTheme.colors.primary else Color.Unspecified
    )

    Row(modifier = modifier.fillMaxHeight(), verticalAlignment = Alignment.CenterVertically) {
        if (startSeparator) {
            Spacer(
                modifier = Modifier
                    .fillMaxHeight(.7f)
                    .width(2.dp)
                    .background(FluentlyTheme.colors.secondary)
            )
        }
        Box(
            modifier = modifier
                .clip(CircleShape)
                .clickable(onClick = onClick)
                .background(color = selectionColor),
        ) {
            Text(
                modifier = Modifier.align(Alignment.Center),
                text = date.dayOfMonth.toString()
            )
        }
        if (endSeparator) {
            Spacer(
                modifier = Modifier
                    .fillMaxHeight(.7f)
                    .width(2.dp)
                    .background(FluentlyTheme.colors.secondary)
            )
        }
    }

}

@DevicePreviews
@Composable
fun CalendarScrollPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface),
            contentAlignment = Alignment.Center
        ) {
            CalendarScroll(
                modifier = Modifier
                    .fillMaxWidth(.8f),
                selectedDay = LocalDate.now().minusDays(10),
                onDayClick = {}
            )
        }
    }
}