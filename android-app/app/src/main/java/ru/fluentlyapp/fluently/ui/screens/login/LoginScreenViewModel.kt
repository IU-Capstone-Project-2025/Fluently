package ru.fluentlyapp.fluently.ui.screens.login

import android.content.Intent
import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import net.openid.appauth.AuthorizationException
import net.openid.appauth.AuthorizationResponse
import ru.fluentlyapp.fluently.oauth.GoogleOAuthService
import javax.inject.Inject

enum class LoginState {
    IDLE,
    SUCCESS,
    ERROR
}

@HiltViewModel
class LoginScreenViewModel @Inject constructor(
    private val googleOAuthService: GoogleOAuthService
) : ViewModel() {
    private val _loginState = MutableStateFlow(LoginState.IDLE)
    val loginState = _loginState.asStateFlow()

    fun getOpenAuthPageIntent() = googleOAuthService.getOpenAuthPageIntent()

    fun handleAuthResponseIntent(dataIntent: Intent?) {
        if (dataIntent == null) {
            _loginState.update { LoginState.ERROR }
            return
        }

        val tokenRequest = AuthorizationResponse.fromIntent(dataIntent)?.createTokenExchangeRequest()
        val exception = AuthorizationException.fromIntent(dataIntent)

        when {
            exception != null -> {
                Log.e("LoginScreenViewModel", "Exception when getting auth code:\n$exception")
                _loginState.update { LoginState.ERROR }
            }
            tokenRequest != null -> {
                // Try to fetch the token from the google
                viewModelScope.launch {
                    try {
                        val token = googleOAuthService.performTokenRequest(tokenRequest)
                        Log.i("LoginScreenViewModel", "Fetched the token: $token; idToken=${token.idToken}")
                        _loginState.update { LoginState.SUCCESS }
                    } catch (ex : Exception) {
                        if (ex is CancellationException) throw ex
                        Log.e("LoginScreenViewModel", "Error when performing token request: $ex")
                        _loginState.update { LoginState.ERROR }
                    }
                }
            }
        }
    }

}