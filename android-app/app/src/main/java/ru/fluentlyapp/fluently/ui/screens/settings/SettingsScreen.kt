package ru.fluentlyapp.fluently.ui.screens.settings

import android.widget.Space
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
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Modifier
import androidx.hilt.navigation.compose.hiltViewModel
import ru.fluentlyapp.fluently.ui.theme.FluentlyTheme
import ru.fluentlyapp.fluently.ui.utils.SmallPhonePreview
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import ru.fluentlyapp.fluently.R
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.ui.theme.components.TopAppBar
import ru.fluentlyapp.fluently.ui.theme.components.UserPreferencesSettings

@Composable
fun SettingsScreen(
    modifier: Modifier = Modifier,
    settingsScreenViewModel: SettingsScreenViewModel = hiltViewModel(),
    onBackClick: () -> Unit,
    onUserLoggedOut: () -> Unit,
) {
    val uiState by settingsScreenViewModel.uiState.collectAsState()
    SettingsScreenContent(
        modifier = modifier,
        uiState = uiState,
        onBackClick = onBackClick,
        onUpdateUserPreferences = {},
        onSubmitUserPreferences = {},
        onUserLoggedOut = {}
    )
    LaunchedEffect(onUserLoggedOut) {
        withContext(Dispatchers.Main.immediate) {
            settingsScreenViewModel.commands.collect { command ->
                if (command is SettingScreenCommand.LoginCredentialsRemovedCommand) {
                    onUserLoggedOut()
                }
            }
        }
    }
}

@Composable
fun SettingsScreenContent(
    modifier: Modifier = Modifier,
    uiState: SettingsScreenUiState,
    onBackClick: () -> Unit,
    onUpdateUserPreferences: (UserPreferences) -> Unit,
    onSubmitUserPreferences: () -> Unit,
    onUserLoggedOut: () -> Unit
) {
    Column(
        modifier = modifier
            .background(FluentlyTheme.colors.surface)
            .windowInsetsPadding(WindowInsets.systemBars)
    ) {
        TopAppBar(
            title = stringResource(R.string.settings),
            modifier = Modifier.fillMaxWidth(),
            onBackClick = onBackClick
        )
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .weight(1f)
                .padding(horizontal = 16.dp)
                .verticalScroll(rememberScrollState())
        ) {
            Spacer(modifier = Modifier.height(8.dp))
            UserPreferencesSettings(
                modifier = Modifier.fillMaxWidth(),
                userPreferences = uiState.userPreferences,
                onPreferencesChange = onUpdateUserPreferences,
                availableTopics = uiState.availableTopics
            )
            Spacer(modifier = Modifier.height(24.dp))
            Column(
                modifier = Modifier.fillMaxWidth(),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Box(
                    modifier = Modifier
                        .clip(RoundedCornerShape(16.dp))
                        .clickable(
                            onClick = onSubmitUserPreferences,
                            enabled = uiState.uploadingState in setOf(
                                SettingsUploading.IDLE, SettingsUploading.ERROR
                            )
                        )
                        .background(FluentlyTheme.colors.primary)
                        .padding(20.dp)
                ) {
                    if (uiState.uploadingState == SettingsUploading.UPLOADING) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(20.dp),
                            color = FluentlyTheme.colors.onPrimary
                        )
                    } else {
                        Text(
                            text = "Update Settings",
                            color = FluentlyTheme.colors.onPrimary,
                            fontSize = 20.sp,
                            fontWeight = FontWeight.Bold
                        )
                    }
                }
                when (uiState.uploadingState) {
                    SettingsUploading.ERROR -> {
                        Text(
                            text = "Error, please try again",
                            color = FluentlyTheme.colors.error
                        )
                    }

                    SettingsUploading.SUCCESS -> {
                        Text(
                            text = "Successfully updated!",
                            color = FluentlyTheme.colors.correct
                        )
                    }

                    else -> {}
                }
            }
            Spacer(modifier = Modifier.height(16.dp))
            Box(
                modifier = Modifier.fillMaxWidth(),
                contentAlignment = Alignment.Center
            ) {
                Row(
                    modifier = Modifier.clickable(onClick = onUserLoggedOut),
                    verticalAlignment = Alignment.CenterVertically,
                ) {
                    Icon(
                        painter = painterResource(R.drawable.ic_logout),
                        tint = FluentlyTheme.colors.onSurfaceVariant,
                        contentDescription = "logout"
                    )
                    Spacer(modifier = Modifier.width(16.dp))
                    Text(
                        text = "Log Out",
                        color = FluentlyTheme.colors.onSurfaceVariant
                    )
                }
            }
        }
    }
}

@SmallPhonePreview
@Composable
fun SettingsScreenPreview() {
    FluentlyTheme {
        Box(modifier = Modifier.fillMaxSize()) {
            SettingsScreenContent(
                modifier = Modifier.fillMaxSize(),
                uiState = SettingsScreenUiState(
                    userPreferences = UserPreferences.empty(),
                    uploadingState = SettingsUploading.UPLOADING,
                    availableTopics = listOf("Travelling", "Animals", "Brain rot")
                ),
                onBackClick = {},
                onUpdateUserPreferences = {},
                onSubmitUserPreferences = {},
                onUserLoggedOut = {}
            )
        }
    }
}