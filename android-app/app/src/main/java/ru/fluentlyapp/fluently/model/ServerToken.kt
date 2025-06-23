package ru.fluentlyapp.fluently.model

data class ServerToken(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresInSeconds: Int
)