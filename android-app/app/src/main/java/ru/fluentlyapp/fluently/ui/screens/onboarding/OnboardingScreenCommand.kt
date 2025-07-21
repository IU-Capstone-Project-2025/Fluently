package ru.fluentlyapp.fluently.ui.screens.onboarding

sealed interface OnboardingScreenCommand {
    object UserPreferencesUploadedCommand : OnboardingScreenCommand
}