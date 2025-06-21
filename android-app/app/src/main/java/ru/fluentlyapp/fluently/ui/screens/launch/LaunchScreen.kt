package ru.fluentlyapp.fluently.ui.screens.launch

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.sp
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHost
import androidx.navigation.NavHostController
import androidx.navigation.compose.rememberNavController
import kotlinx.coroutines.delay
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.navigation.Destination
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import kotlin.time.Duration.Companion.milliseconds

@Composable
fun LaunchScreen(
    modifier: Modifier = Modifier,
    navHostController: NavHostController,
    launchScreenViewModel: LaunchScreenViewModel = viewModel(),
) {
    LaunchScreenContent(modifier = modifier, navHostController)
}

@Composable
fun LaunchScreenContent(
    modifier: Modifier = Modifier,
    navHostController: NavHostController
) {
    LaunchedEffect(Unit) {
        delay(1000.milliseconds)
        navHostController.navigate(Destination.HomeScreen)
    }

    Box(modifier = modifier.background(color = Color.White)) {
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
    LaunchScreenContent(
        modifier = Modifier.fillMaxSize(),
        rememberNavController()
    )
}