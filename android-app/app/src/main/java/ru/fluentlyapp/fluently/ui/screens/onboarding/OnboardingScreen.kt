package ru.fluentlyapp.fluently.ui.screens.onboarding

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.CefrLevel
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.theme.components.MaterialDropDownForm
import ru.fluentlyapp.fluently.ui.theme.components.UserPreferencesSettings
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview

@Composable
fun OnboardingScreen(
    modifier: Modifier = Modifier,
    onComplete: () -> Unit,
    onboardingViewModel: OnboardingViewModel = hiltViewModel()
) {
    val uiState = onboardingViewModel.uiState.collectAsState()
    OnboardingScreenContent(
        modifier = modifier,
        uiState = uiState.value,
        onComplete = {
            onboardingViewModel.completeUserPreferences()
        },
        onUpdaterUserPreferences = {
            onboardingViewModel.updateUserPreferences(it)
        },
        onRetryIfError = { onboardingViewModel.initOnboarding() }
    )

    LaunchedEffect(onComplete) {
        withContext(Dispatchers.Main.immediate) {
            onboardingViewModel.commands.collect { command ->
                if (command is OnboardingScreenCommand.UserPreferencesUploadedCommand) {
                    onComplete()
                }
            }
        }
    }
}

@Composable
fun OnboardingScreenContent(
    modifier: Modifier = Modifier,
    uiState: OnboardingScreenUiState,
    onUpdaterUserPreferences: (UserPreferences) -> Unit,
    onComplete: () -> Unit,
    onRetryIfError: () -> Unit,
) {
    Box(modifier = modifier.background(FluentlyTheme.colors.surface)) {
        when (uiState.initialLoadingState) {
            InitialLoadingState.LOADING -> {
                Column(
                    modifier = Modifier.align(Alignment.Center),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        stringResource(R.string.loading),
                        color = FluentlyTheme.colors.primary,
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    CircularProgressIndicator(color = FluentlyTheme.colors.secondary)
                }
            }

            InitialLoadingState.ERROR -> {
                Column(
                    modifier = Modifier.align(Alignment.Center),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        text = stringResource(R.string.something_went_wrong),
                        color = FluentlyTheme.colors.primary,
                        fontSize = 24.sp,
                        fontWeight = FontWeight.Bold
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    Box(
                        modifier = Modifier
                            .clip(RoundedCornerShape(16.dp))
                            .clickable(onClick = onRetryIfError)
                            .background(FluentlyTheme.colors.secondary)
                            .padding(16.dp)
                    ) {
                        Text(
                            text = stringResource(R.string.retry),
                            color = FluentlyTheme.colors.onSecondary,
                            fontSize = 24.sp
                        )
                    }
                }
            }

            InitialLoadingState.SUCCESS -> {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(16.dp)
                        .verticalScroll(rememberScrollState()),
                    verticalArrangement = Arrangement.Center,
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Text(
                        color = FluentlyTheme.colors.onSurface,
                        text = stringResource(R.string.onboarding_top_text)
                    )
                    Spacer(modifier = Modifier.height(16.dp))
                    UserPreferencesSettings(
                        modifier = Modifier.fillMaxSize(),
                        userPreferences = uiState.userPreferences,
                        onPreferencesChange = onUpdaterUserPreferences,
                        availableTopics = uiState.availableTopics
                    )
                    Spacer(modifier = Modifier.height(32.dp))
                    Box(
                        modifier = Modifier
                            .clip(RoundedCornerShape(16.dp))
                            .clickable(
                                onClick = onComplete,
                                enabled = uiState.uploadingLoadingState in setOf(
                                    UploadingLoadingState.ERROR, UploadingLoadingState.IDLE
                                )
                            )
                            .background(FluentlyTheme.colors.primary)
                            .padding(20.dp)
                    ) {
                        when (uiState.uploadingLoadingState) {
                            UploadingLoadingState.IDLE, UploadingLoadingState.SUCCESS -> {
                                Text(
                                    fontSize = 20.sp,
                                    fontWeight = FontWeight.Bold,
                                    text = stringResource(R.string.finish),
                                    color = FluentlyTheme.colors.onPrimary
                                )
                            }

                            UploadingLoadingState.ERROR -> {
                                Text(
                                    fontSize = 20.sp,
                                    fontWeight = FontWeight.Bold,
                                    text = stringResource(R.string.error_home),
                                    color = FluentlyTheme.colors.onPrimary
                                )
                            }

                            UploadingLoadingState.UPLOADING -> {
                                CircularProgressIndicator(
                                    modifier = Modifier.size(20.dp),
                                    color = FluentlyTheme.colors.onPrimary
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

@SmallPhonePreview
@Composable
fun OnboardingScreenPreview() {
    FluentlyTheme {
        Box(modifier = Modifier.fillMaxSize()) {
            OnboardingScreenContent(
                modifier = Modifier.fillMaxSize(),
                uiState = OnboardingScreenUiState(
                    initialLoadingState = InitialLoadingState.SUCCESS,
                    userPreferences = UserPreferences(
                        avatarImageUrl = "",
                        cefrLevel = CefrLevel.A1, // Choose this
                        factEveryday = false,
                        goal = "", // Choose this
                        id = "",
                        notifications = true,
                        subscribed = false,
                        userId = "",
                        wordsPerDay = 10 // Choose this
                    ),
                    availableTopics = listOf("travel", "work", "life"),
                    uploadingLoadingState = UploadingLoadingState.UPLOADING
                ),
                onComplete = {},
                onRetryIfError = {},
                onUpdaterUserPreferences = {}
            )
        }
    }
}
