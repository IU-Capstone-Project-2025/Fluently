package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.model.ServerToken

interface FluentlyNetworkDataSource {
    suspend fun getServerToken(idToken: String): ServerToken
    suspend fun refreshServerToken(refreshToken: String): ServerToken
}