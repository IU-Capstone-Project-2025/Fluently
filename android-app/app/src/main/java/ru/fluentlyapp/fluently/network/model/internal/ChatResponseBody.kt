package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class ChatResponseBody(
    val chat: List<MessageApiModel>
)
