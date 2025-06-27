package ru.fluentlyapp.fluently.ui.screens.lesson.components.other

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun LoadingLessonComponent(
    modifier: Modifier
) {
    Column(
        modifier = modifier.background(color = FluentlyTheme.colors.surface),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = stringResource(R.string.lesson_is_loading),
            fontWeight = FontWeight.Bold,
            fontSize = 24.sp,
            color = FluentlyTheme.colors.secondary
        )
        Spacer(modifier = Modifier.height(16.dp))
        CircularProgressIndicator(
            color = FluentlyTheme.colors.primary
        )
    }
}

@Preview(device = PIXEL_7)
@Composable
fun LoadingLessonComponentPreview() {
    FluentlyTheme {
        LoadingLessonComponent(modifier = Modifier.fillMaxSize())
    }
}
