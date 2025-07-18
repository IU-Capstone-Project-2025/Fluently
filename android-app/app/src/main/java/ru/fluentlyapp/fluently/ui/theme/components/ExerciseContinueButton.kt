package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun ExerciseContinueButton(
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    onClick: () -> Unit
) {
    if (enabled) {
        Box(
            modifier = modifier
                .clip(RoundedCornerShape(16.dp))
                .clickable(onClick = onClick, enabled = enabled)
                .background(color = FluentlyTheme.colors.primary)
                .padding(20.dp)
        ) {
            Text(
                modifier = Modifier.align(Alignment.Center),
                text = stringResource(R.string.continue_to_next_exercise),
                fontSize = 20.sp,
                fontWeight = FontWeight.Bold,
                color = FluentlyTheme.colors.onPrimary
            )
        }
    }
    else {
        Box(
            modifier = modifier
                .padding(bottom = 32.dp)
                .fillMaxWidth()
        )
    }
}

@Composable
@Preview
fun ExerciseContinueButtonPreview() {
    FluentlyTheme {
        ExerciseContinueButton(
            enabled = true,
            onClick = {}
        )
    }
}
