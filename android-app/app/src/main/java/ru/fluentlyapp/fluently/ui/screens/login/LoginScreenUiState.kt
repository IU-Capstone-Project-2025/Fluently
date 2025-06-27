package ru.fluentlyapp.fluently.ui.screens.login

enum class LoginLoadingState {
    IDLE,
    LOADING,
    ERROR,
    SUCCESS
}

data class LoginScreenUiState(
    val loginLoadingState: LoginLoadingState
)