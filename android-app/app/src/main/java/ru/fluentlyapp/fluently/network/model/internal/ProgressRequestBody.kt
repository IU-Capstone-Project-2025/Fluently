package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class WordProgressApiModel(
    val cnt_reviewed: Int? = null,
    val confidence_score: Int? = null,
    val learned_at: String? = null,
    val word_id: String
)