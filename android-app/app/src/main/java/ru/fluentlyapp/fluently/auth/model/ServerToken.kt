package ru.fluentlyapp.fluently.auth.model

import kotlinx.serialization.Serializable

@Serializable
data class ServerToken(
    val accessToken: String,
    val refreshToken: String,
    val tokenType: String,
    val expiresInSeconds: Int
)