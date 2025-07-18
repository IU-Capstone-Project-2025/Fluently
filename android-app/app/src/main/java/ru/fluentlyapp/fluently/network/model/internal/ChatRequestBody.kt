package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class ChatRequestBody(
    val chat: List<MessageApiModel>
)

@Serializable
data class MessageApiModel(
    val author: String,
    val message: String
)