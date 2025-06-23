package ru.fluentlyapp.fluently.data.repository

import android.content.Context
import android.content.Intent
import dagger.hilt.android.qualifiers.ApplicationContext
import net.openid.appauth.TokenRequest
import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.network.FluentlyNetworkDataSource
import ru.fluentlyapp.fluently.oauth.GoogleOAuthService
import ru.fluentlyapp.fluently.oauth.model.OAuthToken
import javax.inject.Inject

interface AuthRepository {
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
     * Replace the currently saved `ServerToken` with the passed one.
     */
    suspend fun updateServerToken(serverToken: ServerToken)

    /**
     * Try to refresh the token using the currently saved `ServerToken`.
     * The methods DOES NOT save the fetched `ServerToken`.
     *
     * In case of `ServerToken` isn't stored, throws exception. May throw other exceptions.
     */
    suspend fun refreshServerToken(): ServerToken

    /**
     * Delete (if saved) the `ServerToken` locally
     */
    suspend fun deleteServerToken()
}


class GoogleBasedAuthRepository @Inject constructor(
    @ApplicationContext private val applicationContext: Context,
    private val googleOAuthService: GoogleOAuthService,
    private val fluentlyNetworkDataSource: FluentlyNetworkDataSource
) : AuthRepository {
    override suspend fun isUserLogged(): Boolean {
        TODO("Not yet implemented")
    }

    override fun getAuthPageIntent(): Intent {
        return googleOAuthService.getOpenAuthPageIntent()
    }

    override suspend fun getOAuthToken(tokenRequest: TokenRequest): OAuthToken {
        return googleOAuthService.performTokenRequest(tokenRequest)
    }

    override suspend fun getServerToken(oauthToken: OAuthToken): ServerToken {
        return fluentlyNetworkDataSource.getServerToken(oauthToken.idToken)
    }

    override suspend fun updateServerToken(serverToken: ServerToken) {
        TODO("Not yet implemented")
    }

    override suspend fun refreshServerToken(): ServerToken {
        TODO("Not yet implemented")
    }

    override suspend fun deleteServerToken() {
        TODO("Not yet implemented")
    }

}