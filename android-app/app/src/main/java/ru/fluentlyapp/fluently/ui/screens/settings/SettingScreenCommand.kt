package ru.fluentlyapp.fluently.ui.screens.settings

sealed interface SettingScreenCommand {
    object LoginCredentialsRemovedCommand : SettingScreenCommand
    object SettingsUpdatedCommand : SettingScreenCommand
}