package ru.fluentlyapp.fluently.ui.screens.lesson.components.decoration

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

abstract class LearningPartCompleteObserver {
    abstract fun onMoveNext()
}

@Composable
fun LearningPartCompleteDecoration(
    modifier: Modifier = Modifier,
    learningPartCompleteObserver: LearningPartCompleteObserver
) {
    Column(
        modifier = modifier
            .background(color = FluentlyTheme.colors.surface)
            .padding(horizontal = 16.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            textAlign = TextAlign.Center,
            text = stringResource(R.string.learning_part_complete_text),
            fontSize = 24.sp,
            lineHeight = 32.sp
        )
        Spacer(modifier = Modifier.height(32.dp))
        Box(
            modifier = Modifier
                .clip(RoundedCornerShape(16.dp))
                .clickable(onClick = { learningPartCompleteObserver.onMoveNext() })
                .background(color = FluentlyTheme.colors.primary)
                .padding(20.dp)
        ) {
            Text(
                text = "Continue",
                color = FluentlyTheme.colors.onPrimary,
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold
            )
        }
    }
}

@Preview(device = PIXEL_7)
@Composable
fun LearningPartCompleteDecorationPreview() {
    FluentlyTheme {
        LearningPartCompleteDecoration(
            modifier = Modifier.fillMaxSize(),
            learningPartCompleteObserver = object : LearningPartCompleteObserver() {
                override fun onMoveNext() {}
            }
        )
    }
}
