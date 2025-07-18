package ru.fluentlyapp.fluently.feature.joinedwordprogress

import java.time.Instant

data class JoinedWordProgress(
    val id: String,
    val word: String,
    val translation: String,
    val examples: List<Pair<String, String>>,
    val isLearning: Boolean,
    val instant: Instant
)