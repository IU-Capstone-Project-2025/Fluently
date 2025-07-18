package ru.fluentlyapp.fluently.common.model

import kotlinx.serialization.Serializable

@Serializable
data class UserPreferences(
    val avatarImageUrl: String,
    val cefrLevel: String,
    val factEveryday: Boolean,
    val goal: String,
    val id: String,
    val notifications: Boolean,
    val subscribed: Boolean,
    val userId: String,
    val wordsPerDay: Int
)