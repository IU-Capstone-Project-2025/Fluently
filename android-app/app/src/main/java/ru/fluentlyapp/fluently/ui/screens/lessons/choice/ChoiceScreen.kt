package ru.fluentlyapp.fluently.ui.screens.lessons.choice

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.selects.select
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.ThemeColors

@Composable
fun ChoiceScreenContent(
    modifier: Modifier = Modifier,
    uiState: ChoiceScreenUiState,
    onContinueClick: () -> Unit,
    onVariantClick: (Int) -> Unit,
) {
    Column(
        modifier = modifier.background(color = Color.White),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f)
                .verticalScroll(
                    state = rememberScrollState()
                ),
            verticalArrangement = Arrangement.Bottom,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(uiState.word, fontSize = 40.sp)

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                stringResource(R.string.choose_correct_translation),
                fontSize = 16.sp
            )

            Spacer(modifier = Modifier.height(8.dp))

            val correctAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = Color.Green,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val wrongAnswerModifier = remember {
                Modifier.border(
                    width = 2.dp,
                    color = Color.Red,
                    shape = RoundedCornerShape(12.dp)
                )
            }

            val notChosenModifier = remember {
                Modifier.alpha(alpha = .4f)
            }

            repeat(times = uiState.variants.size) { index ->
                val itemModifier = when {
                    uiState.selectedVariant == null -> Modifier.clickable(
                        onClick = { onVariantClick(index) },
                    )

                    uiState.correctVariant == index -> correctAnswerModifier
                    uiState.selectedVariant == index -> wrongAnswerModifier
                    else -> notChosenModifier
                }

                Text(
                    modifier = Modifier
                            then itemModifier
                        .clip(shape = RoundedCornerShape(12.dp))
                        .fillMaxWidth(fraction = 0.6f)
                        .background(color = Color.LightGray)
                        .padding(16.dp),
                    text = uiState.variants[index],
                )
                if (index != uiState.variants.size - 1) {
                    Spacer(modifier = Modifier.height(4.dp))
                }
            }
        }

        Spacer(modifier = Modifier.height(64.dp))

        Box(modifier = Modifier.weight(.5f)) {
            androidx.compose.animation.AnimatedVisibility(
                visible = uiState.selectedVariant != null
            ) {
                Box(
                    modifier = Modifier
                        .clickable(onClick = onContinueClick)
                        .clip(shape = RoundedCornerShape(12.dp))
                        .fillMaxWidth(.6f)
                        .background(color = ThemeColors.primary)
                        .padding(16.dp)
                ) {
                    Text(
                        text = stringResource(R.string.continue_to_next_exercise),
                        fontSize = 24.sp,
                        textAlign = TextAlign.Center,
                        modifier = Modifier.align(Alignment.Center),
                        color = Color.White
                    )
                }
            }
        }
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun ChoiceScreenPreview() {
    ChoiceScreenContent(
        modifier = Modifier.fillMaxSize(),
        uiState = ChoiceScreenUiState(
            word = "Influence",
            variants = listOf(
                "Влияние", "Благодарнсость", "Двойственность", "Комар"
            ),
            correctVariant = 0,
            selectedVariant = 2,
        ),
        onContinueClick = {},
        onVariantClick = {}
    )
}