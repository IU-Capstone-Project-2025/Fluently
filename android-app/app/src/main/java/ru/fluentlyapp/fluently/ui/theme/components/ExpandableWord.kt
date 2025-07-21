package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.core.animateFloatAsState
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
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.rotate
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

@Composable
fun ExpandableWord(
    modifier: Modifier = Modifier,
    word: WordUiState,
    isExpanded: Boolean,
    onClick: () -> Unit
) {

    Column(
        modifier = modifier
            .clip(RoundedCornerShape(size = 10.dp))
            .clickable(onClick = onClick)
            .background(color = FluentlyTheme.colors.surfaceContainerHigh)
            .padding(10.dp),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.Start
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(
                    text = word.word,
                    fontWeight = FontWeight.SemiBold,
                    fontSize = 20.sp,
                    color = FluentlyTheme.colors.onSurface
                )
                Text(
                    text = word.translation,
                    color = FluentlyTheme.colors.onSurface,
                    fontSize = 14.sp
                )
            }

            val rotationAngle by animateFloatAsState(
                targetValue = if (!isExpanded) 0f else 90f
            )
            Icon(
                painter = painterResource(R.drawable.ic_chevron_right),
                contentDescription = "Expand list",
                modifier = Modifier
                    .rotate(rotationAngle)
                    .size(32.dp)
            )
        }
        AnimatedVisibility(visible = isExpanded) {
            Column(modifier = Modifier.fillMaxWidth()) {
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = stringResource(R.string.Examples),
                    color = FluentlyTheme.colors.onSurfaceVariant
                )
                repeat(word.examples.size) { index ->
                    val (english, translation) = word.examples[index]
                    Text(text = english)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(text = translation)
                    if (index != word.examples.size - 1) {
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
                                    .border(
                                        width = 1.dp,
                                        color = FluentlyTheme.colors.onSurfaceVariant
                                    )
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
@DevicePreviews
fun ExpandableWordPreview() {
    FluentlyTheme {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(state = rememberScrollState())
                .background(FluentlyTheme.colors.surface),
            contentAlignment = Alignment.TopCenter
        ) {
            var isExpanded by remember { mutableStateOf(false) }
            ExpandableWord(
                modifier = Modifier.fillMaxWidth(.8f),
                word = WordUiState(
                    word = "summer",
                    translation = "лагерь",
                    examples = listOf(
                        "I left Bobby Jr. at an amusement park. Told his father he was at summer camp." to
                                "Я Бобби как-то забыла в луна-парке, а отцу сказала, что он в лагере, выгадала неделю.",
                        "In summer camp the kids used to pay me to make their beds." to
                                "В лагере дети платили мне, чтобы я убирала кровати.",
                    )
                ),
                isExpanded = isExpanded,
                onClick = { isExpanded = !isExpanded }
            )
        }
    }
}
