package ru.fluentlyapp.fluently.auth.oauth

import android.content.Context
import android.content.Intent
import androidx.browser.customtabs.CustomTabsIntent
import androidx.core.net.toUri
import dagger.hilt.android.qualifiers.ApplicationContext
import net.openid.appauth.AuthorizationRequest
import net.openid.appauth.AuthorizationService
import net.openid.appauth.AuthorizationServiceConfiguration
import net.openid.appauth.TokenRequest
import okhttp3.OkHttpClient
import ru.fluentlyapp.fluently.auth.model.OAuthToken
import ru.fluentlyapp.fluently.common.di.BaseOkHttpClient
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton
import kotlin.coroutines.suspendCoroutine

@Singleton
class GoogleOAuthService @Inject constructor(
    @ApplicationContext applicationContext: Context
) {
    private val authService = AuthorizationService(applicationContext)

    private val serviceConfiguration = AuthorizationServiceConfiguration(
        GoogleOAuthConfig.AUTH_URI.toUri(),
        GoogleOAuthConfig.TOKEN_URI.toUri(),
    )

    private fun getAuthRequest(): AuthorizationRequest {
        return AuthorizationRequest.Builder(
            serviceConfiguration,
            GoogleOAuthConfig.CLIENT_ID,
            GoogleOAuthConfig.RESPONSE_TYPE,
            GoogleOAuthConfig.REDIRECT_URI.toUri()
        )
            .setScope(GoogleOAuthConfig.SCOPE)
            .build()
    }

    fun getOpenAuthPageIntent(): Intent {
        val customTabsIntent = CustomTabsIntent.Builder().build()

        val openAuthPageIntent = authService.getAuthorizationRequestIntent(
            getAuthRequest(),
            customTabsIntent
        )
        return openAuthPageIntent
    }

    /**
     * Perform a request for OAuthToken using the passed `tokenRequest`
     *
     * Throws exception if the token retrieval failed
     */
    suspend fun performTokenRequest(
        tokenRequest: TokenRequest
    ): OAuthToken {
        return suspendCoroutine { continuation ->
            authService.performTokenRequest(tokenRequest) { response, ex ->
                when {
                    // The token retrieval is successful
                    response != null -> {
                        // We expect all tokens to be present in the response
                        val token = OAuthToken(
                            accessToken = response.accessToken!!,
                            refreshToken = response.refreshToken!!,
                            idToken = response.idToken!!
                        )
                        continuation.resumeWith(Result.success(token))
                    }
                    // The token retrieval failed -> throw exception
                    ex != null -> {
                        continuation.resumeWith(Result.failure(ex))
                    }
                }
            }
        }
    }
}