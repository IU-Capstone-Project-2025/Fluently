package ru.fluentlyapp.fluently.network.middleware

import android.util.Log
import kotlinx.coroutines.runBlocking
import okhttp3.Authenticator
import okhttp3.Request
import okhttp3.Response
import okhttp3.Route
import ru.fluentlyapp.fluently.auth.AuthManager
import ru.fluentlyapp.fluently.auth.model.ServerToken
import ru.fluentlyapp.fluently.auth.datastore.ServerTokenDataStore
import ru.fluentlyapp.fluently.network.HEADER_AUTHORIZATION
import ru.fluentlyapp.fluently.network.TOKEN_TYPE
import ru.fluentlyapp.fluently.auth.api.ServerTokenApiService
import javax.inject.Inject
import javax.inject.Singleton

/**
 * If the server returned 401 code, this interceptor will be called
 *
 * The interceptor will try to get the refresh token using the `RefreshServerTokenService`
 */
@Singleton
class AuthAuthenticator @Inject constructor(
    private val serverTokenDataStore: ServerTokenDataStore,
    private val refreshServerTokenService: ServerTokenApiService,
    private val authManager: AuthManager
) : Authenticator {
    override fun authenticate(route: Route?, response: Response): Request? {
        val token = runBlocking {
            authManager.getSavedServerToken()
        }
        synchronized(this) {
            val updatedToken = runBlocking {
                authManager.getSavedServerToken()
            }
            val accessToken: String? = if (updatedToken != token) {
                // While the thread was blocked on the synchronize block, some other thread
                // has already fetched the fresh token
                updatedToken?.accessToken
            } else {
                // Otherwise, fetch and store the refresh token
                val newServerToken: ServerToken? = runBlocking {
                    try {
                        authManager.sendRefreshToken()
                    } catch (ex: Exception) {
                        Log.e("Authentificator", "Couldn't fetch the server token in refresh send: $ex")
                        null
                    }
                }

                if (newServerToken != null) {
                    // Save the received server token
                    runBlocking {
                        serverTokenDataStore.saveServerToken(newServerToken)
                    }

                    newServerToken.accessToken
                } else {
                    null
                }
            }

            return if (accessToken != null) {
                response.request
                    .newBuilder()
                    .addHeader(HEADER_AUTHORIZATION, "$TOKEN_TYPE $accessToken")
                    .build()
            } else {
                null
            }
        }
    }
}