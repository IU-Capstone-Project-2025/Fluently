package ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration

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
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

abstract class OnboardingDecorationObserver {
    abstract fun onContinue()
}

data class OnboardingDecorationUiState(
    val newWordsCount: Int,
    val exercisesCount: Int
)

@Composable
fun OnboardingDecoration(
    modifier: Modifier = Modifier,
    onboardingDecorationUiState: OnboardingDecorationUiState,
    onboardingDecorationObserver: OnboardingDecorationObserver,
) {
    Column(
        modifier = modifier
            .background(color = FluentlyTheme.colors.surface)
            .padding(16.dp),
        verticalArrangement = Arrangement.SpaceEvenly,
    ) {
        Column(
            modifier = Modifier.fillMaxWidth(),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Text(
                textAlign = TextAlign.Center,
                text = "Начнём урок!",
                fontSize = 64.sp,
                lineHeight = 64.sp
            )
            Spacer(modifier = Modifier.height(16.dp))
            Text(
                text = "В этом уроке ты:",
                fontSize = 20.sp
            )
            Spacer(modifier = Modifier.height(8.dp))
            // TODO: add normal plural formatting for russian and english
            Text(
                "Узнаешь %d новых слов".format(onboardingDecorationUiState.newWordsCount),
                fontSize = 20.sp
            )
            Spacer(modifier = Modifier.height(4.dp))
            Text(
                "Пройдёшь %d упражнений".format(onboardingDecorationUiState.exercisesCount),
                fontSize = 20.sp
            )
        }

        Box(modifier = Modifier.fillMaxWidth(), contentAlignment = Alignment.Center) {
            Box(
                modifier = Modifier
                    .size(160.dp)
                    .clip(CircleShape)
                    .clickable(onClick = onboardingDecorationObserver::onContinue)
                    .background(FluentlyTheme.colors.primary),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    "Go!",
                    fontWeight = FontWeight.Bold,
                    fontSize = 64.sp,
                    color = FluentlyTheme.colors.onPrimary
                )
            }
        }
    }
}

@Composable
@DevicePreviews
fun OnboardingDecorationPreview() {
    FluentlyTheme {
        OnboardingDecoration(
            modifier = Modifier.fillMaxSize(),
            onboardingDecorationUiState = OnboardingDecorationUiState(10, 78),
            onboardingDecorationObserver = object : OnboardingDecorationObserver() {
                override fun onContinue() {}
            }
        )
    }
}
