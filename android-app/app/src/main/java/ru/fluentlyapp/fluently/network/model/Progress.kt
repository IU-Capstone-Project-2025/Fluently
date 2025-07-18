package ru.fluentlyapp.fluently.network.model

import java.time.Instant

data class Progress(
    val progresses: List<SentWordProgress>
)

data class SentWordProgress(
    val wordId: String,
    val cntReviewed: Int,
    val confidenceScore: Int,
    val learnedAt: Instant
)