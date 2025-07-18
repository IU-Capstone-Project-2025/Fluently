package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class WordProgressApiModel(
    val cnt_reviewed: Int,
    val confidence_score: Int,
    val learned_at: String,
    val word_id: String
)