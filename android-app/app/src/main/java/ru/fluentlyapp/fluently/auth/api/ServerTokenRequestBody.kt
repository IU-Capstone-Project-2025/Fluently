package ru.fluentlyapp.fluently.auth.api

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class GetServerTokenRequestBody(
    @SerialName("id_token")
    val idToken: String,

    val platform: String
)