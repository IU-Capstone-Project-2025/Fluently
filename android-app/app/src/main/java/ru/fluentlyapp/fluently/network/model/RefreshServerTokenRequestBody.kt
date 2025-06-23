package ru.fluentlyapp.fluently.network.model

import kotlinx.serialization.SerialName

data class RefreshServerTokenRequestBody(
    @SerialName("refresh_token")
    val refreshToken: String
)