package ru.fluentlyapp.fluently.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.tooling.preview.Devices.PIXEL_7
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme

@Composable
fun TopAppBar(
    modifier: Modifier = Modifier,
    onBackClick: () -> Unit
) {
    Box(
        modifier = modifier
        .background(color = FluentlyTheme.colors.surface)
        .padding(8.dp)
    ) {
        Icon(
            modifier = Modifier.clickable(
                onClick = onBackClick
            ),
            painter = painterResource(R.drawable.ic_chevron_left),
            contentDescription = "Back button"
        )
    }
}

@Composable
@Preview(device = PIXEL_7)
fun TopBarPreview() {
    FluentlyTheme {
        Box(modifier = Modifier.fillMaxSize()) {
            TopAppBar(
                modifier = Modifier.fillMaxWidth(),
                onBackClick = {}
            )
        }
    }
}