package ru.fluentlyapp.fluently.network.model

data class ServerTokenResponseBody(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresInSeconds: Int
)