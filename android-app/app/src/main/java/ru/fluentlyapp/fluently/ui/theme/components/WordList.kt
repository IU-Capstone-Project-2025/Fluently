package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateSetOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

@Composable
fun WordList(
    modifier: Modifier = Modifier,
    words: List<WordUiState>
) {
    val expandedWords = remember { mutableStateSetOf<Int>() }
    LazyColumn(
        modifier = modifier.padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        items(words.size) { index ->
            ExpandableWord(
                modifier = Modifier.fillMaxWidth(),
                word = words[index],
                isExpanded = expandedWords.contains(index),
                onClick = {
                    if (expandedWords.contains(index))
                        expandedWords.remove(index)
                    else
                        expandedWords.add(index)
                }
            )
        }
    }
}

@Composable
@DevicePreviews
fun WordListPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .background(FluentlyTheme.colors.surface),
            contentAlignment = Alignment.Center
        ) {
            WordList(
                modifier = Modifier.fillMaxWidth(.8f),
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
        }
    }
}