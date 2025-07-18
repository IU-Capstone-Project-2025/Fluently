package ru.fluentlyapp.fluently.feature.wordcache

import kotlinx.serialization.Serializable

@Serializable
data class WordCache(
    val wordId: String,
    val word: String,
    val translation: String,
    val examples: List<Pair<String, String>> // (sentence, translation of the sentence)
)