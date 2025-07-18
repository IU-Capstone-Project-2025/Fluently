package ru.fluentlyapp.fluently.network.model

import kotlinx.serialization.Serializable

@Serializable
data class WordOfTheDay(
    val wordId: String,
    val word: String,
    val translation: String,
    val examples: List<Pair<String, String>>,
)