package ru.fluentlyapp.fluently.feature.wordprogress

import ru.fluentlyapp.fluently.database.app.wordprogress.WordProgressEntity

fun WordProgressEntity.toWordProgress() = WordProgress(
    wordId = id,
    isLearning = isLearning,
    instant = timestamp
)

fun WordProgress.toWordProgressEntity() = WordProgressEntity(
    id = wordId,
    isLearning = isLearning,
    timestamp = instant
)