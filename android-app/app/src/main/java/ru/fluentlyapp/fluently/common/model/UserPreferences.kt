package ru.fluentlyapp.fluently.common.model

import kotlinx.serialization.Serializable

enum class CefrLevel(val key: String) {
    A1("A1"),
    A2("A2"),
    B1("B1"),
    B2("B2"),
    C1("C1"),
    C2("C2")
}

@Serializable
data class UserPreferences(
    val avatarImageUrl: String,
    val cefrLevel: CefrLevel,
    val factEveryday: Boolean,
    val goal: String,
    val id: String,
    val notifications: Boolean,
    val subscribed: Boolean,
    val userId: String,
    val wordsPerDay: Int
)

