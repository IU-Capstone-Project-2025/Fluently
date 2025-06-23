package ru.fluentlyapp.fluently.datastore.model

import kotlinx.serialization.Serializable

@Serializable
data class ServerTokenPreference(
    val accessToken: String,
    val refreshToken: String,
    val expiresInSeconds: Int,
    val tokenType: String
)