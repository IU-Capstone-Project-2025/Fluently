package ru.fluentlyapp.fluently.feature.wordprogress

import java.time.Instant

data class WordProgress(
    val wordId: String,
    val isLearning: Boolean,
    val timestamp: Instant
)