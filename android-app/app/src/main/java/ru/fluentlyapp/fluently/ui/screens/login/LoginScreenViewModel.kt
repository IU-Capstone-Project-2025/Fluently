package ru.fluentlyapp.fluently.ui.screens.login

import android.content.Intent
import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import ru.fluentlyapp.fluently.auth.AuthManager
import javax.inject.Inject


@HiltViewModel
class LoginScreenViewModel @Inject constructor(
    private val authManager: AuthManager
) : ViewModel() {
    private val _uiState = MutableStateFlow(
        LoginScreenUiState(
            loginLoadingState = LoginLoadingState.IDLE
        )
    )
    val uiState = _uiState.asStateFlow()

    fun getOpenAuthPageIntent() = authManager.getAuthPageIntent()

    fun handleAuthResponseIntent(dataIntent: Intent?) {
        if (dataIntent == null) {
            _uiState.update { it.copy(loginLoadingState = LoginLoadingState.ERROR) }
            return
        }


        _uiState.update {
            it.copy(loginLoadingState = LoginLoadingState.LOADING)
        }

        viewModelScope.launch(Dispatchers.IO) {
            try {
                authManager.handleReturnedDataIntent(dataIntent)
                _uiState.update { it.copy(loginLoadingState = LoginLoadingState.SUCCESS) }
            } catch (ex: Exception) {
                Log.e("LoginScreenViewModel", "Exception while handing dataIntent: $ex")
                _uiState.update { it.copy(loginLoadingState = LoginLoadingState.ERROR) }
            }
        }
    }

}