package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class UserPreferencesRequestBody(
    val avatar_image_url: String? = null,
    val cefr_level: String? = null,
    val fact_everyday: Boolean? = null,
    val goal: String? = null,
    val id: String? = null,
    val notifications: Boolean? = null,
    val notification_at: String? = null,
    val subscribed: Boolean? = null,
    val user_id: String? = null,
    val words_per_day: Int? = null
)