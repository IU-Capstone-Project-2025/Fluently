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
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.DevicePreviews

abstract class FinishDecorationObserver {
    abstract fun onFinish()
}

@Composable
fun FinishDecoration(
    modifier: Modifier = Modifier,
    finishDecorationObserver: FinishDecorationObserver
) {
    Column(
        modifier = modifier
            .background(color = FluentlyTheme.colors.surface)
            .padding(16.dp),
        verticalArrangement = Arrangement.SpaceEvenly,
        horizontalAlignment = Alignment.Start
    ) {
        Text(
            textAlign = TextAlign.Center,
            text = "Урок пройден!",
            fontSize = 64.sp,
            lineHeight = 64.sp
        )

        Box(modifier = Modifier.fillMaxWidth(), contentAlignment = Alignment.Center) {
            Box(
                modifier = Modifier
                    .clip(RoundedCornerShape(75.dp))
                    .background(FluentlyTheme.colors.primary)
                    .clickable(onClick = finishDecorationObserver::onFinish)
                    .padding(25.dp),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    "Я молодец \uD83D\uDE0E",
                    fontWeight = FontWeight.Bold,
                    fontSize = 32.sp,
                    color = FluentlyTheme.colors.onPrimary
                )
            }
        }
    }
}

@Composable
@DevicePreviews
fun FinishDecorationPreview() {
    FluentlyTheme {
        FinishDecoration(
            modifier = Modifier.fillMaxSize(),
            finishDecorationObserver = object : FinishDecorationObserver() {
                override fun onFinish() {}
            }
        )
    }
}
