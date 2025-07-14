package ru.fluentlyapp.fluently.feature.joinedwordprogress

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

interface JoinedWordProgressRepository {
    fun getJoinedWordProgresses(
        begin: Instant,
        end: Instant
    ): Flow<List<JoinedWordProgress>>

    fun getAllJoinedWordProgresses(): Flow<List<JoinedWordProgress>>
}


@Singleton
class JoinedWordProgressRepositoryImpl @Inject constructor(
    appDatabase: AppDatabase
) : JoinedWordProgressRepository {
    private val joinedWordProgressDao = appDatabase.joinedWordProgressDao()

    override fun getJoinedWordProgresses(
        begin: Instant,
        end: Instant
    ): Flow<List<JoinedWordProgress>> {
        return joinedWordProgressDao.getJoinedWordProgress(begin, end).map { list ->
            list.map { it.toJoinedWordProgress() }
        }
    }

    override fun getAllJoinedWordProgresses(): Flow<List<JoinedWordProgress>> {
        return joinedWordProgressDao.getAllJoinedWordProgress().map { list ->
            list.map { it.toJoinedWordProgress() }
        }
    }
}