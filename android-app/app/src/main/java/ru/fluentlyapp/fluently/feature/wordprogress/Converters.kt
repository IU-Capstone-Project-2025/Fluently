package ru.fluentlyapp.fluently.feature.wordprogress

import ru.fluentlyapp.fluently.feature.wordprogress.database.WordProgressEntity

fun WordProgress.toWordProgressEntity() = WordProgressEntity(
    id = wordId,
    isLearning = isLearning,
    timestamp = timestamp
)

fun WordProgressEntity.toWordProgress() = WordProgress(
    wordId = id,
    isLearning = isLearning,
    timestamp = timestamp
)