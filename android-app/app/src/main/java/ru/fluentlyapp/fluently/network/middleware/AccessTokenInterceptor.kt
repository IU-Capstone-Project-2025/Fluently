package ru.fluentlyapp.fluently.network.middleware

import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import ru.fluentlyapp.fluently.auth.AuthManager
import ru.fluentlyapp.fluently.auth.model.ServerToken
import ru.fluentlyapp.fluently.auth.datastore.ServerTokenDataStore
import ru.fluentlyapp.fluently.network.HEADER_AUTHORIZATION
import ru.fluentlyapp.fluently.network.TOKEN_TYPE
import javax.inject.Inject
import javax.inject.Singleton

/**
 * Appends the access token received from the `ServerTokenManager` to the
 * If the access token is not stored it just adds null
 */
@Singleton
class AccessTokenInterceptor @Inject constructor(
    private val authManager: AuthManager
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token: ServerToken? = runBlocking {
            authManager.getSavedServerToken()
        }

        val updatedRequest = chain
            .request()
            .newBuilder()
            .addHeader(HEADER_AUTHORIZATION, "$TOKEN_TYPE $token")
            .build()
        return chain.proceed(updatedRequest)
    }
}