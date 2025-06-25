package ru.fluentlyapp.fluently.navigation

import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreen
import ru.fluentlyapp.fluently.ui.screens.launch.LaunchScreen
import ru.fluentlyapp.fluently.ui.screens.lesson.LessonFlowScreen
import ru.fluentlyapp.fluently.ui.screens.login.LoginScreen

@Composable
fun FluentlyNavHost(
    modifier: Modifier = Modifier,
    navHostController: NavHostController = rememberNavController()
) {
    NavHost(
        modifier = modifier,
        navController = navHostController,
        startDestination = Destination.HomeScreen
    ) {
        composable<Destination.LaunchScreen> {
            LaunchScreen(
                modifier = Modifier.fillMaxSize(),
                onUserLogged = {
                    navHostController.navigate(Destination.HomeScreen)
                },
                onUserNotLogged = {
                    navHostController.navigate(Destination.LoginScreen)
                }
            )
        }

        composable<Destination.LoginScreen> {
            LoginScreen(
                modifier = Modifier.fillMaxSize(),
                onSuccessfulLogin = {
                    navHostController.navigate(Destination.HomeScreen)
                }
            )
        }

        composable<Destination.HomeScreen> {
            HomeScreen(
                modifier = Modifier.fillMaxSize(),
                onNavigateToLesson = {
                    navHostController.navigate(Destination.LessonScreen("lesson_id"))
                }
            )
        }

        composable<Destination.LessonScreen> {
            LessonFlowScreen(
                modifier = Modifier.fillMaxSize(),
                onBackClick = {
                    navHostController.navigate(Destination.HomeScreen)
                }
            )
        }
    }
}