package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.data.model.ServerToken

interface FluentlyNetworkDataSource {
    /**
     * Get the ServerToken from the passed idToken. The idToken is the token received from
     * the OAuthFlow after submitting the authorization code.
     *
     * The method may throw an exception in case of network or any other type of error.
     */
    suspend fun getServerToken(idToken: String): ServerToken

    /**
     * Get the ServerToken from the passed refreshToken. The refreshToken is one of the fields
     * of the ServerToken class.
     *
     * The method may throw an exception in case of network or any other type of error.
     */
    suspend fun refreshServerToken(refreshToken: String): ServerToken
}