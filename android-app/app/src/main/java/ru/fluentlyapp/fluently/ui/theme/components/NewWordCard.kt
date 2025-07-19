package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxHeight
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

data class WordUiState(
    val word: String,
    val translation: String,
    val examples: List<Pair<String, String>>
)

@Composable
fun NewWordCard(
    modifier: Modifier,
    word: String,
    translation: String,
    examples: List<Pair<String, String>> // (sentence, translation of that sentence)
) {
    Column(
        modifier = modifier
            .clip(RoundedCornerShape(size = 16.dp))
            .background(color = FluentlyTheme.colors.surfaceContainerHigh)
            .padding(16.dp),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Column(
            modifier = Modifier
                .verticalScroll(state = rememberScrollState(), overscrollEffect = null)
        ) {
            Text(
                text = word,
                fontSize = 32.sp,
                color = FluentlyTheme.colors.onSurface
            )
            Spacer(modifier = Modifier.height(16.dp))
            Text(
                text = stringResource(R.string.translation),
                color = FluentlyTheme.colors.onSurfaceVariant
            )
            Text(
                text = translation
            )
            Spacer(modifier = Modifier.height(16.dp))

            Text(
                text = stringResource(R.string.examples),
                color = FluentlyTheme.colors.onSurfaceVariant
            )
            repeat(examples.size) { index ->
                val (english, translation) = examples[index]
                Text(text = english)
                Spacer(modifier = Modifier.height(8.dp))
                Text(text = translation)
                if (index != examples.size - 1) {
                    Box(
                        modifier = Modifier
                            .height(16.dp)
                            .fillMaxWidth(),
                        contentAlignment = Alignment.Center
                    ) {
                        Spacer(
                            modifier = Modifier
                                .fillMaxWidth()
                                .height(1.dp)
                                .border(width = 1.dp, color = FluentlyTheme.colors.onSurfaceVariant)
                        )
                    }
                }
            }
        }
    }
}

@Composable
@DevicePreviews
fun NewWordCardPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .background(FluentlyTheme.colors.surface)
                .fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            NewWordCard(
                modifier = Modifier
                    .fillMaxWidth(fraction = .8f)
                    .fillMaxHeight(fraction = .7f),
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
    }
}

@Composable
@DevicePreviews
fun NewWordCardScrollPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .background(FluentlyTheme.colors.surface)
                .fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            NewWordCard(
                modifier = Modifier
                    .fillMaxWidth(fraction = .8f)
                    .fillMaxHeight(fraction = .7f),
                word = "summer",
                translation = "лагерь",
                examples = listOf(
                    "I left Bobby Jr. at an amusement park. Told his father he was at summer camp." to
                            "Я Бобби как-то забыла в луна-парке, а отцу сказала, что он в лагере, выгадала неделю.",
                    "In summer camp the kids used to pay me to make their beds." to
                            "В лагере дети платили мне, чтобы я убирала кровати.",
                    "We need a nice, clean towel here at summer camp." to
                            "Нам нужно милое, и чистое полотенце в этом лагере.",
                    "We need a nice, clean towel here at summer camp." to
                            "Нам нужно милое, и чистое полотенце в этом лагере.",
                    "We need a nice, clean towel here at summer camp." to
                            "Нам нужно милое, и чистое полотенце в этом лагере.",

                    )
            )
        }
    }
}