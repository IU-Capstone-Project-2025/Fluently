package ru.fluentlyapp.fluently.data.repository

import android.content.Context
import android.content.Intent
import dagger.hilt.android.qualifiers.ApplicationContext
import net.openid.appauth.TokenRequest
import retrofit2.HttpException
import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.datastore.ServerTokenDataStore
import ru.fluentlyapp.fluently.network.model.GetServerTokenRequestBody
import ru.fluentlyapp.fluently.network.model.RefreshServerTokenRequest
import ru.fluentlyapp.fluently.network.services.ServerTokenApiService
import ru.fluentlyapp.fluently.network.toServerToken
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
    private val googleOAuthService: GoogleOAuthService,
    private val serverTokenService: ServerTokenApiService,
    private val serverTokenDataStore: ServerTokenDataStore
) : AuthRepository {

    override suspend fun isUserLogged(): Boolean {
        return serverTokenDataStore.getServerToken() != null
    }

    override fun getAuthPageIntent(): Intent {
        return googleOAuthService.getOpenAuthPageIntent()
    }

    override suspend fun getOAuthToken(tokenRequest: TokenRequest): OAuthToken {
        return googleOAuthService.performTokenRequest(tokenRequest)
    }

    override suspend fun getServerToken(oauthToken: OAuthToken): ServerToken {
        return serverTokenService.getServerToken(
            GetServerTokenRequestBody(
                idToken = oauthToken.idToken,
                platform = "android"
            )
        ).body()?.toServerToken()!!
    }

    override suspend fun updateServerToken(serverToken: ServerToken) {
        serverTokenDataStore.saveServerToken(serverToken)
    }

    override suspend fun refreshServerToken(): ServerToken {
        val serverToken = serverTokenDataStore.getServerToken()

        return serverTokenService.refreshToken(
            RefreshServerTokenRequest(
                refreshToken = serverToken!!.refreshToken
            )
        ).body()?.toServerToken()!!
    }

    override suspend fun deleteServerToken() {
        serverTokenDataStore.deleteServerToken()
    }
}