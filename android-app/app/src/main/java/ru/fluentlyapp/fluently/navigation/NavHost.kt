package ru.fluentlyapp.fluently.navigation

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import ru.fluentlyapp.fluently.ui.screens.launch.LaunchScreen
import ru.fluentlyapp.fluently.ui.screens.login.LoginScreen

@Composable
fun FluentlyNavHost(
    modifier: Modifier = Modifier,
    navHostController: NavHostController = rememberNavController()
) {
    NavHost(
        modifier = modifier,
        navController = navHostController,
        startDestination = Destination.LaunchScreen
    ) {
        composable<Destination.LaunchScreen> {
            LaunchScreen(
                modifier = Modifier.fillMaxSize(),
                navHostController
            )
        }

        composable<Destination.LoginScreen> {
            LoginScreen(
                modifier = Modifier.fillMaxSize()
            )
        }

    }
}