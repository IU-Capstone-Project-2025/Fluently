package ru.fluentlyapp.fluently.ui.screens.login

import android.util.Log
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.animation.AnimatedContent
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.expandVertically
import androidx.compose.animation.shrinkVertically
import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.systemBars
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.windowInsetsPadding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.R
import timber.log.Timber

@Composable
fun LoginScreen(
    modifier: Modifier = Modifier,
    loginScreenViewModel: LoginScreenViewModel = hiltViewModel(),
    onSuccessfulLogin: () -> Unit,
) {
    val uiState by loginScreenViewModel.uiState.collectAsStateWithLifecycle()

    if (uiState.loginLoadingState == LoginLoadingState.SUCCESS) {
        onSuccessfulLogin()
    }

    val authPageLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { activityResult ->
        Timber.d("Received activityResult=$activityResult")
        loginScreenViewModel.handleAuthResponseIntent(activityResult.data)
    }

    LoginScreenContent(
        modifier = modifier,
        onLoginWithGoogleClick = {
            authPageLauncher.launch(loginScreenViewModel.getOpenAuthPageIntent())
        },
        uiState = uiState
    )
}

@Composable
fun LoginScreenContent(
    modifier: Modifier = Modifier,
    onLoginWithGoogleClick: () -> Unit,
    uiState: LoginScreenUiState
) {
    Column(
        modifier = modifier
            .background(color = FluentlyTheme.colors.primary)
            .windowInsetsPadding(WindowInsets.systemBars)
    ) {
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .height(240.dp)
        ) {
            Text(
                modifier = Modifier
                    .align(Alignment.Center)
                    .fillMaxWidth(),
                textAlign = TextAlign.Center,
                lineHeight = 48.sp,
                text = stringResource(R.string.welcome_to_fluently),
                fontSize = 48.sp,
                color = FluentlyTheme.colors.onPrimary,
                fontWeight = FontWeight.Bold
            )
        }
        Column(
            modifier = Modifier
                .clip(RoundedCornerShape(topStart = 32.dp, topEnd = 32.dp))
                .fillMaxWidth()
                .weight(1f)
                .background(color = FluentlyTheme.colors.surface),
            verticalArrangement = Arrangement.Center,
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Row(
                modifier = Modifier
                    .clip(shape = RoundedCornerShape(8.dp))
                    .clickable(onClick = onLoginWithGoogleClick)
                    .background(color = FluentlyTheme.colors.secondary)
                    .padding(top = 4.dp, bottom = 4.dp, start = 4.dp, end = 16.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                AnimatedContent(
                    modifier = Modifier
                        .size(40.dp),
                    targetState = uiState.loginLoadingState
                ) { targetLoadingState ->
                    when {
                        targetLoadingState == LoginLoadingState.LOADING -> {
                            CircularProgressIndicator(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .padding(4.dp),
                                color = FluentlyTheme.colors.onSecondary
                            )
                        }

                        else -> {
                            Icon(
                                modifier = Modifier
                                    .clip(shape = RoundedCornerShape(8.dp))
                                    .fillMaxSize()
                                    .background(color = FluentlyTheme.colors.onSecondary)
                                    .padding(4.dp),
                                painter = painterResource(R.drawable.google_logo),
                                contentDescription = null,
                                tint = Color.Unspecified
                            )
                        }
                    }
                }
                Spacer(modifier = Modifier.width(16.dp))
                val authText = if (uiState.loginLoadingState == LoginLoadingState.LOADING) {
                    stringResource(R.string.entering_text)
                } else {
                    stringResource(R.string.sign_in_with_google_account)
                }
                Text(
                    text = authText,
                    fontSize = 16.sp,
                    color = FluentlyTheme.colors.onSecondary
                )
            }

            AnimatedVisibility(
                visible = uiState.loginLoadingState == LoginLoadingState.ERROR,
                enter = expandVertically(),
                exit = shrinkVertically()
            ) {
                Column {
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(horizontal = 8.dp),
                        textAlign = TextAlign.Center,
                        fontWeight = FontWeight.SemiBold,
                        color = FluentlyTheme.colors.error,
                        text = stringResource(R.string.login_error_text)
                    )
                }
            }
        }
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun LoginScreenPreview() {
    FluentlyTheme {
        var testLoadingState by remember {
            mutableStateOf(LoginLoadingState.ERROR)
        }

        LoginScreenContent(
            modifier = Modifier.fillMaxSize(),
            onLoginWithGoogleClick = {
                val entryIndex = (testLoadingState.ordinal + 1) % LoginLoadingState.entries.size
                testLoadingState =
                    LoginLoadingState.entries[entryIndex]
            },
            uiState = LoginScreenUiState(
                loginLoadingState = testLoadingState
            )
        )
    }
}