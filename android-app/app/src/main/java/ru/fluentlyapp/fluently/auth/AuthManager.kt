package ru.fluentlyapp.fluently.auth

import android.content.Intent
import android.net.Uri
import androidx.core.net.toUri
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.supervisorScope
import kotlinx.coroutines.withContext
import net.openid.appauth.AuthorizationException
import net.openid.appauth.AuthorizationResponse
import net.openid.appauth.TokenRequest
import okhttp3.OkHttpClient
import okhttp3.Request
import ru.fluentlyapp.fluently.auth.model.ServerToken
import ru.fluentlyapp.fluently.auth.datastore.ServerTokenDataStore
import ru.fluentlyapp.fluently.auth.api.GetServerTokenRequestBody
import ru.fluentlyapp.fluently.auth.api.RefreshServerTokenRequest
import ru.fluentlyapp.fluently.auth.api.ServerTokenApiService
import ru.fluentlyapp.fluently.auth.api.toServerToken
import ru.fluentlyapp.fluently.auth.oauth.GoogleOAuthService
import ru.fluentlyapp.fluently.auth.model.OAuthToken
import ru.fluentlyapp.fluently.common.di.BaseOkHttpClient
import ru.fluentlyapp.fluently.datastore.UserPreferencesDataStore
import timber.log.Timber
import javax.inject.Inject

interface AuthManager {
    /**
     * Returns `true` if the app have saved credentials `false` otherwise
     */
    suspend fun isUserLogged(): Boolean

    /**
     * Get the intent that opens the custom tab intent that shows the page fetched from the
     * authorization server (usually google api).
     */
    fun getAuthPageIntent(): Intent

    /**
     * After receiving the ActivityResult from the auth page, pass the `dataIntent` to this method
     * The method is convenience wrapper for all steps required to obtain and save `ServerToken`
     *
     * May throw exception.
     */
    suspend fun handleReturnedDataIntent(dataIntent: Intent)

    /**
     * Send the `tokenRequest` to the authorization server and fetches the OAuthToken.
     *
     * May throw exception.
     */
    suspend fun getOAuthToken(tokenRequest: TokenRequest): OAuthToken

    /**
     * Get the serverToken from the passed OAuthToken.
     *
     * May throw exception.
     */
    suspend fun getServerToken(oauthToken: OAuthToken): ServerToken

    /**
     * Return the saved `ServerToken`.
     *
     * Returns null is the `ServerToken` is not stored.
     */
    suspend fun getSavedServerToken(): ServerToken?

    /**
     * Replace the currently saved `ServerToken` with the passed one.
     */
    suspend fun updateServerToken(serverToken: ServerToken)

    /**
     * Send the field `refreshToken` of the stored `ServerToken` to exchange it for a new
     * `ServerToken`.
     *
     * In case of `ServerToken` isn't stored, throws exception. May throw other exceptions.
     */
    suspend fun sendRefreshToken(): ServerToken

    /**
     * Delete (if saved) the `ServerToken` locally
     */
    suspend fun deleteServerToken()
}

class GoogleBasedOAuthManager @Inject constructor(
    @BaseOkHttpClient private val baseOkHttpClient: OkHttpClient,
    private val userPreferencesDataStore: UserPreferencesDataStore,
    private val googleOAuthService: GoogleOAuthService,
    private val serverTokenApiService: ServerTokenApiService,
    private val serverTokenDataStore: ServerTokenDataStore,
) : AuthManager {

    override suspend fun isUserLogged(): Boolean {
        return serverTokenDataStore.getServerToken() != null
    }

    override fun getAuthPageIntent(): Intent {
        return googleOAuthService.getOpenAuthPageIntent()
    }

    override suspend fun handleReturnedDataIntent(dataIntent: Intent) {
        val tokenRequest = AuthorizationResponse.fromIntent(dataIntent)?.createTokenExchangeRequest()
        val exception = AuthorizationException.fromIntent(dataIntent)

        when {
            exception != null -> {
                throw exception
            }
            tokenRequest != null -> {
                // Try to fetch the token from the google
                handleTokenRequest(tokenRequest)
            }
        }
    }

    private suspend fun handleTokenRequest(tokenRequest: TokenRequest) {
        val token = getOAuthToken(tokenRequest)
        Timber.d("Fetch the OAuth Token: $token")

        val serverToken = getServerToken(token)
        Timber.d("Fetch the server token from the OAuth token: $serverToken")

        updateServerToken(serverToken)
        Timber.d("Successfully save the server token")
    }

    override suspend fun getOAuthToken(tokenRequest: TokenRequest): OAuthToken {
        return googleOAuthService.performTokenRequest(tokenRequest)
    }

    override suspend fun getServerToken(oauthToken: OAuthToken): ServerToken {
        return serverTokenApiService.getServerToken(
            GetServerTokenRequestBody(
                idToken = oauthToken.idToken,
                platform = "android"
            )
        ).body()?.toServerToken()!!
    }

    override suspend fun updateServerToken(serverToken: ServerToken) {
        serverTokenDataStore.saveServerToken(serverToken)
    }

    override suspend fun getSavedServerToken(): ServerToken? {
        return serverTokenDataStore.getServerToken()
    }

    override suspend fun sendRefreshToken(): ServerToken {
        val serverToken = serverTokenDataStore.getServerToken()

        return serverTokenApiService.refreshToken(
            RefreshServerTokenRequest(
                refreshToken = serverToken!!.refreshToken
            )
        ).body()?.toServerToken()!!
    }

    override suspend fun deleteServerToken() {
        serverTokenDataStore.deleteServerToken()
    }
}