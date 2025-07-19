package ru.fluentlyapp.fluently.ui.screens.launch

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.systemBars
import androidx.compose.foundation.layout.windowInsetsPadding
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.delay
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import kotlin.time.Duration.Companion.milliseconds
import androidx.compose.runtime.getValue
import androidx.hilt.navigation.compose.hiltViewModel

@Composable
fun LaunchScreen(
    modifier: Modifier = Modifier,
    launchScreenViewModel: LaunchScreenViewModel = hiltViewModel(),
    onUserLogged: () -> Unit,
    onUserNotLogged: () -> Unit
) {
    val isUserLogged by launchScreenViewModel.isUserLogged.collectAsState()

    LaunchedEffect(isUserLogged) {
        if (isUserLogged != null) {
            delay(700.milliseconds) // Just wait some time to avoid flickering screen

            if (isUserLogged == true) {
                onUserLogged()
            } else if (isUserLogged == false) {
                onUserNotLogged()
            }
        }
    }

    LaunchScreenContent(modifier = modifier)
}

@Composable
fun LaunchScreenContent(
    modifier: Modifier = Modifier,
) {
    Box(
        modifier = modifier
            .background(color = Color.White)
            .windowInsetsPadding(WindowInsets.systemBars)
    ) {
        Text(
            modifier = Modifier.align(Alignment.Center),
            text = stringResource(R.string.fluently),
            fontSize = 64.sp,
            fontWeight = FontWeight.Bold,
            color = FluentlyTheme.colors.primary
        )
    }
}

@Composable
@Preview(device = Devices.PIXEL_7)
fun LaunchScreenPreview() {
    FluentlyTheme {
        LaunchScreenContent(
            modifier = Modifier.fillMaxSize(),
        )
    }
}