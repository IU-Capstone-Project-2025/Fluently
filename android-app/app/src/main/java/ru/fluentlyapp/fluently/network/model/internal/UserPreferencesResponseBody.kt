package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class UserPreferencesResponseBody(
    val avatar_image_url: String,
    val cefr_level: String,
    val fact_everyday: Boolean,
    val goal: String,
    val id: String,
    val notifications: Boolean,
    val subscribed: Boolean,
    val user_id: String,
    val words_per_day: Int
)