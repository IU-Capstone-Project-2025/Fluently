package ru.fluentlyapp.fluently.ui.screens.wordsprogress

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.theme.components.WordList
import ru.fluentlyapp.fluently.ui.theme.components.WordUiState
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

@Composable
fun WordsProgressScreen(
    modifier: Modifier = Modifier,
    onBackClick: () -> Unit,
    wordsInProgressViewModel: WordsInProgressViewModel = hiltViewModel()
) {
    val uiState by wordsInProgressViewModel.uiState.collectAsStateWithLifecycle()
    WordsProgressScreenContent(
        modifier = modifier,
        uiState = uiState,
        onBackClick = onBackClick,
    )
}

@Composable
fun WordsProgressScreenContent(
    modifier: Modifier = Modifier,
    uiState: WordsProgressUiState,
    onBackClick: () -> Unit,
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
                text = uiState.pageTitle,
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
                .padding(top = 32.dp, start = 20.dp, end = 20.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.SpaceEvenly
        ) {
            SearchBar(
                modifier = Modifier.fillMaxWidth(),
                onQueryChange = {},
                onSearch = {},
                query = ""
            )
            Spacer(modifier = Modifier.height(16.dp))
            WordList(
                modifier = Modifier.fillMaxSize(),
                words = uiState.words.map {
                    WordUiState(
                        word = it.word,
                        translation = it.translation,
                        examples = it.examples
                    )
                }
            )
        }
    }
}

@DevicePreviews
@Composable
fun WordsProgressScreenPreview() {
    FluentlyTheme {
        WordsProgressScreenContent(
            modifier = Modifier.fillMaxSize(),
            onBackClick = {},
            uiState = WordsProgressUiState(
                pageTitle = "In progress",
                searchString = "aboba",
                words = (1..100).map {
                    WordUiState(
                        word = "summer",
                        translation = "лагерь",
                        examples = listOf(
                            "I left Bobby Jr. at an amusement park. Told his father he was at summer camp." to
                                    "Я Бобби как-то забыла в луна-парке, а отцу сказала, что он в лагере, выгадала неделю.",
                            "In summer camp the kids used to pay me to make their beds." to
                                    "В лагере дети платили мне, чтобы я убирала кровати.",
                        )
                    )
                }
            )
        )
    }
}