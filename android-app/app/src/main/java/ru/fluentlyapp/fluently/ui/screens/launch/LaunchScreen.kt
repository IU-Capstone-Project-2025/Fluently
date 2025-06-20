package ru.fluentlyapp.fluently.ui.screens.launch

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.lifecycle.viewmodel.compose.viewModel
import ru.fluentlyapp.fluently.R

@Composable
fun LaunchScreen(
    modifier: Modifier,
    launchScreenViewModel: LaunchScreenViewModel = viewModel(),
) {
    LaunchScreenContent(modifier = modifier)
}

@Composable
fun LaunchScreenContent(
    modifier: Modifier,
) {
    Box(modifier = Modifier.fillMaxSize()) {
        Text(
            text = stringResource(R.string.fluently),
            style = MaterialTheme.typography.headlineMedium
        )
    }
}