package ru.fluentlyapp.fluently.feature.joinedwordprogress

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class JoinedWordProgressRepository @Inject constructor(
    appDatabase: AppDatabase
) {
    private val joinedWordProgressDao = appDatabase.joinedWordProgressDao()

    fun getJoinedWordProgresses(
        begin: Instant,
        end: Instant
    ): Flow<List<JoinedWordProgress>> {
        return joinedWordProgressDao.getJoinedWordProgress(begin, end).map { list ->
            list.map {
                val wordCacheEntity = Json.decodeFromString<WordCache>(it.wordJson)
                JoinedWordProgress(
                    id = it.id,
                    word = wordCacheEntity.word,
                    translation = wordCacheEntity.translation,
                    examples = wordCacheEntity.examples,
                    isLearning = it.isLearning,
                    instant = it.timestamp
                )
            }
        }
    }
}