package ru.fluentlyapp.fluently.network.model

import kotlinx.serialization.Serializable

@Serializable
data class ServerTokenResponseBody(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresInSeconds: Int
)