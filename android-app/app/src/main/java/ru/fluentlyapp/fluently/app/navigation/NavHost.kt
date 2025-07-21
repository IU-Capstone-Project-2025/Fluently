package ru.fluentlyapp.fluently.app.navigation

import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.navigation.NavHostController
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import ru.fluentlyapp.fluently.ui.screens.calendar.CalendarScreen
import ru.fluentlyapp.fluently.ui.screens.home.HomeScreen
import ru.fluentlyapp.fluently.ui.screens.launch.LaunchScreen
import ru.fluentlyapp.fluently.ui.screens.lesson.LessonFlowScreen
import ru.fluentlyapp.fluently.ui.screens.login.LoginScreen
import ru.fluentlyapp.fluently.ui.screens.onboarding.OnboardingScreen
import ru.fluentlyapp.fluently.ui.screens.settings.SettingsScreen
import ru.fluentlyapp.fluently.ui.screens.wordsprogress.WordsProgressScreen

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
                onUserLogged = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.LaunchScreen> {
                            inclusive = true
                        }
                    }
                },
                onUserNotLogged = {
                    navHostController.navigate(Destination.LoginScreen) {
                        popUpTo<Destination.LaunchScreen> {
                            inclusive = true
                        }
                    }
                }
            )
        }

        composable<Destination.LoginScreen> {
            LoginScreen(
                modifier = Modifier.fillMaxSize(),
                onSuccessfulLogin = {
                    navHostController.navigate(Destination.OnboardingScreen) {
                        popUpTo<Destination.LoginScreen>() {
                            inclusive = true
                        }
                    }
                }
            )
        }

        composable<Destination.HomeScreen> {
            HomeScreen(
                modifier = Modifier.fillMaxSize(),
                onNavigateToLesson = {
                    navHostController.navigate(Destination.LessonScreen)
                },
                onNavigateToCalendar = {
                    navHostController.navigate(Destination.CalendarScreen)
                },
                onLearnedWordsClick = {
                    navHostController.navigate(
                        Destination.WordsProgress(isLearning = false)
                    )
                },
                onInProgressWordsClick = {
                    navHostController.navigate(
                        Destination.WordsProgress(isLearning = true)
                    )
                },
                onNavigateToSettings = {
                    navHostController.navigate(Destination.SettingsScreen)
                }
            )
        }

        composable<Destination.LessonScreen> {
            LessonFlowScreen(
                modifier = Modifier.fillMaxSize(),
                onBackClick = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.HomeScreen>()
                        launchSingleTop = true
                    }
                }
            )
        }

        composable<Destination.WordsProgress> {
            WordsProgressScreen(
                modifier = Modifier.fillMaxSize(),
                onBackClick = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.HomeScreen>()
                        launchSingleTop = true
                    }
                }
            )
        }

        composable<Destination.CalendarScreen> {
            CalendarScreen(
                modifier = Modifier.fillMaxSize(),
                onBackClick = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.HomeScreen>()
                        launchSingleTop = true
                    }
                },
            )
        }

        composable<Destination.OnboardingScreen> {
            OnboardingScreen(
                modifier = Modifier.fillMaxSize(),
                onComplete = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.LoginScreen>() {
                            inclusive = true
                        }
                    }
                },
            )
        }

        composable<Destination.SettingsScreen> {
            SettingsScreen(
                modifier = Modifier.fillMaxSize(),
                onBackClick = {
                    navHostController.navigate(Destination.HomeScreen) {
                        popUpTo<Destination.HomeScreen>()
                        launchSingleTop = true
                    }
                },
                onUserLoggedOut = {
                    navHostController.navigate(Destination.LaunchScreen) {
                        popUpTo<Destination.HomeScreen> {
                            inclusive = true
                        }
                    }
                }
            )
        }
    }
}