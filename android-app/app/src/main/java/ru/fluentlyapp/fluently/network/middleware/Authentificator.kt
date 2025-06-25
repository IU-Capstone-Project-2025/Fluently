package ru.fluentlyapp.fluently.network.middleware

import kotlinx.coroutines.runBlocking
import okhttp3.Authenticator
import okhttp3.Request
import okhttp3.Response
import okhttp3.Route
import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.datastore.ServerTokenDataStore
import ru.fluentlyapp.fluently.network.HEADER_AUTHORIZATION
import ru.fluentlyapp.fluently.network.TOKEN_TYPE
import ru.fluentlyapp.fluently.network.model.RefreshServerTokenRequest
import ru.fluentlyapp.fluently.network.services.ServerTokenApiService
import ru.fluentlyapp.fluently.network.toServerToken
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
    private val refreshServerTokenService: ServerTokenApiService
) : Authenticator {
    override fun authenticate(route: Route?, response: Response): Request? {
        val token = runBlocking {
            serverTokenDataStore.getServerToken()
        }
        synchronized(this) {
            val updatedToken = runBlocking {
                serverTokenDataStore.getServerToken()
            }
            val accessToken: String? = if (updatedToken != token) {
                // While the thread was blocked on the synchronize block, some other thread
                // has already fetched the fresh token
                updatedToken?.accessToken
            } else {
                // Otherwise, fetch and store the refresh token
                val currentServerToken: ServerToken? = runBlocking {
                    serverTokenDataStore.getServerToken()
                }

                if (currentServerToken == null) {
                    // The server token isn't stored
                    return null
                }

                // Using the current `ServerToken` and the stored refresh token, fetch the new `ServerToken`
                val refreshServiceResponse = runBlocking {
                    refreshServerTokenService.refreshToken(
                        RefreshServerTokenRequest(refreshToken = currentServerToken.refreshToken)
                    )
                }

                val refreshServiceResponseBody = refreshServiceResponse.body()

                if (!refreshServiceResponse.isSuccessful || refreshServiceResponseBody == null) {
                    // For some reason, the service has failed
                    return null
                }

                val newServerToken: ServerToken = refreshServiceResponseBody.toServerToken()

                // Save the new server token
                runBlocking {
                    serverTokenDataStore.saveServerToken(newServerToken)
                }

                newServerToken.accessToken
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