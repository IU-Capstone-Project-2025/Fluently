package ru.fluentlyapp.fluently.oauth.model

data class OAuthToken(
    val accessToken: String,
    val refreshToken: String,
    val idToken: String
)