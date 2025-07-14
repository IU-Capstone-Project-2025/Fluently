package ru.fluentlyapp.fluently.feature.joinedwordprogress

import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.database.app.joinedwordprogress.JoinedWordProgressData
import ru.fluentlyapp.fluently.feature.wordcache.WordCache

fun JoinedWordProgressData.toJoinedWordProgress(): JoinedWordProgress {
    val wordCacheEntity = Json.decodeFromString<WordCache>(wordJson)
    return JoinedWordProgress(
        id = id,
        word = wordCacheEntity.word,
        translation = wordCacheEntity.translation,
        examples = wordCacheEntity.examples,
        isLearning = isLearning,
        instant = timestamp
    )
}