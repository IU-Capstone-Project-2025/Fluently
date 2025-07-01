package ru.fluentlyapp.fluently.auth.api

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class RefreshServerTokenRequest(
    @SerialName("refresh_token")
    val refreshToken: String
)