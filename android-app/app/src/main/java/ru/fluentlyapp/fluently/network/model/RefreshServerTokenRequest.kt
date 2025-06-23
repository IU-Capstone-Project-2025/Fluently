package ru.fluentlyapp.fluently.network.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class RefreshServerTokenRequest(
    @SerialName("refresh_token")
    val refreshToken: String
)