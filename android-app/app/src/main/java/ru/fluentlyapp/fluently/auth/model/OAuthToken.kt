package ru.fluentlyapp.fluently.auth.model

data class OAuthToken(
    val accessToken: String,
    val refreshToken: String,
    val idToken: String
)