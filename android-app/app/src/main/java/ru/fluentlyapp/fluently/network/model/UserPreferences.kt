package ru.fluentlyapp.fluently.network.model

data class UserPreferences(
    val avatarImageUrl: String,
    val cefrLevel: String,
    val factEveryday: Boolean,
    val goal: String,
    val id: String,
    val notificationAt: String,
    val notifications: Boolean,
    val subscribed: Boolean,
    val userId: String,
    val wordsPerDay: Int
)