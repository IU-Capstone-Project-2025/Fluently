package ru.fluentlyapp.fluently.data.model

data class ServerToken(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresInSeconds: Int
)