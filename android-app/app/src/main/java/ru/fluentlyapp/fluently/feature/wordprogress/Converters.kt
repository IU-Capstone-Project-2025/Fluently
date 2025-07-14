package ru.fluentlyapp.fluently.feature.wordprogress

import ru.fluentlyapp.fluently.database.app.wordprogress.WordProgressEntity

fun WordProgressEntity.toWordProgress() = WordProgress(
    wordId = wordId,
    isLearning = isLearning,
    instant = timestamp
)

fun WordProgress.toWordProgressEntity() = WordProgressEntity(
    wordId = wordId,
    isLearning = isLearning,
    timestamp = instant
)