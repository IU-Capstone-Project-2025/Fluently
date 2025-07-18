package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun MessageBubble(
    modifier: Modifier = Modifier,
    text: String,
    fromUser: Boolean
) {
    val messageBgColor = if (fromUser) {
        FluentlyTheme.colors.primary
    } else {
        FluentlyTheme.colors.primaryVariant
    }

    val messageFont = if (fromUser) {
        FluentlyTheme.colors.onPrimary
    } else {
        FluentlyTheme.colors.onSurface
    }
    Box(
        modifier = Modifier
            .clip(RoundedCornerShape(16.dp))
            .background(color = messageBgColor)
            .padding(horizontal = 16.dp, vertical = 8.dp)
    ) {
        Text(
            text = text,
            color = messageFont
        )
    }
}