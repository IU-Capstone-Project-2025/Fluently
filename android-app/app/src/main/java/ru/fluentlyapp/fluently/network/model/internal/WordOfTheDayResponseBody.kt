package ru.fluentlyapp.fluently.network.model.internal

import kotlinx.serialization.Serializable

@Serializable
data class WordOfTheDayResponseBody(
    val word_id: String,
    val word: String,
    val translation: String,
    val is_learned: Boolean,
    val topic: String,
    val sentences: List<SentenceApiModel>,
)
