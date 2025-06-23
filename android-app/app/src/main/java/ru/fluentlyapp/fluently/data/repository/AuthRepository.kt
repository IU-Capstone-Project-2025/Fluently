package ru.fluentlyapp.fluently.data.repository

import android.content.Intent
import net.openid.appauth.TokenRequest
import ru.fluentlyapp.fluently.oauth.model.OAuthToken

interface AuthRepository {
    /**
     * Returns `true` if the app have saved credentials `false` otherwise
     * It possible for credentials to be outdated (i.e. access_token has expired)
     * In this case use `isServerTokenFresh()` to check the freshness
     */
    suspend fun isUserLogged(): Boolean

    /**
     * Get the intent that opens the custom tab intent that shows the page fetched from the
     * authorization server (usually google api)
     */
    fun getAuthPageIntent(): Intent

    /**
     * Sends the `tokenRequest` to the authorization server and fetches the OAuthToken
     * May throw exception in case of internet error or malformed tokenRequest
     */
    suspend fun getOAuthToken(tokenRequest: TokenRequest): OAuthToken

    /**
     *
     */
    suspend fun getServerToken(oauthToken: OAuthToken)
    suspend fun updateServerToken()
    suspend fun isServerTokenFresh(): Boolean
    suspend fun refreshServerToken()
}