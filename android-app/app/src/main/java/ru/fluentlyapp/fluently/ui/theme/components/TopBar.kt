package ru.fluentlyapp.fluently.ui.theme.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.statusBars
import androidx.compose.foundation.layout.systemBars
import androidx.compose.foundation.layout.windowInsetsPadding
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun TopAppBar(
    modifier: Modifier = Modifier,
    title: String? = null,
    onBackClick: () -> Unit,
) {
    Box(
        modifier = modifier
            .background(color = FluentlyTheme.colors.surface)
            .windowInsetsPadding(WindowInsets.statusBars)
            .padding(8.dp),
    ) {
        Icon(
            modifier = Modifier
                .clip(CircleShape)
                .clickable(onClick = onBackClick),
            painter = painterResource(R.drawable.ic_chevron_left),
            contentDescription = "Back button"
        )
        if (title != null) {
            Text(modifier = Modifier.align(Alignment.Center), text = title, fontSize = 20.sp)
        }
    }
}

@Composable
@Preview(device = PIXEL_7)
fun TopBarPreview() {
    FluentlyTheme {
        Box(modifier = Modifier.fillMaxSize()) {
            TopAppBar(
                modifier = Modifier.fillMaxWidth(),
                onBackClick = {},
                title = "Settings"
            )
        }
    }
}