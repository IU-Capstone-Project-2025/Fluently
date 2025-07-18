package ru.fluentlyapp.fluently.feature.wordcache

import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.database.app.wordcache.WordCacheEntity

fun WordCacheEntity.toWordCache(): WordCache {
    return Json.decodeFromString<WordCache>(wordJson)
}

fun WordCache.toWordCacheEntity(): WordCacheEntity {
    return WordCacheEntity(
        id = wordId,
        wordJson = Json.encodeToString(this)
    )
}