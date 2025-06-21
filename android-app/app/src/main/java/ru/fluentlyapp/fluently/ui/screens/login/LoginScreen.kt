package ru.fluentlyapp.fluently.ui.screens.login

import android.util.Log
import androidx.activity.compose.rememberLauncherForActivityResult
import androidx.activity.result.contract.ActivityResultContracts
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Devices
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.compose.runtime.getValue

@Composable
fun LoginScreen(
    modifier: Modifier = Modifier,
    loginScreenViewModel: LoginScreenViewModel = hiltViewModel()
) {
    val loginState by loginScreenViewModel.loginState.collectAsStateWithLifecycle()

    val authPageLauncher = rememberLauncherForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { activityResult ->
        Log.i("LoginScreen", "Received activityResult=$activityResult")
        loginScreenViewModel.handleAuthResponseIntent(activityResult.data)
    }

    LoginScreenContent(
        modifier = modifier,
        onLoginWithGoogleClick = {
            authPageLauncher.launch(loginScreenViewModel.getOpenAuthPageIntent())
        },
        loginState = loginState
    )
}

@Composable
fun LoginScreenContent(
    modifier: Modifier = Modifier,
    onLoginWithGoogleClick: () -> Unit,
    loginState: LoginState
) {
    Column(
        modifier = modifier.background(color = Color.White),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text("Login to Fluently", fontSize = 32.sp)

        Spacer(modifier = Modifier.height(16.dp))

        Button(
            contentPadding = PaddingValues(16.dp),
            onClick = onLoginWithGoogleClick
        ) {
            Text("continue with Google", fontSize = 20.sp)
        }

        when (loginState) {
            LoginState.SUCCESS -> {
                Spacer(modifier = Modifier.height(16.dp))
                Text("Success!", color = Color.Green)
            }
            LoginState.ERROR -> {
                Spacer(modifier = Modifier.height(16.dp))
                Text("Error :(", color = Color.Red)
            }
            LoginState.IDLE -> {
                // Do not produce anything
            }
        }
    }
}

@Preview(device = Devices.PIXEL_7)
@Composable
fun LoginScreenPreview() {
    LoginScreenContent(
        modifier = Modifier.fillMaxSize(),
        {},
        LoginState.ERROR
    )
}