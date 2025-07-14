package ru.fluentlyapp.fluently.feature.joinedwordprogress

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import timber.log.Timber
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

interface JoinedWordProgressRepository {
    fun getJoinedWordProgresses(
        begin: Instant,
        end: Instant
    ): Flow<List<JoinedWordProgress>>

    fun getPerWordOverallProgress(): Flow<List<JoinedWordProgress>>
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

    override fun getPerWordOverallProgress(): Flow<List<JoinedWordProgress>> {
        return joinedWordProgressDao.getAllJoinedWordProgress().map { list ->
            Timber.d("getPerWordOverallProgress; tmp result: $list")
            val map = mutableMapOf<String, JoinedWordProgress>()
            list.forEach {
                if (!map.contains(it.id) || (map[it.id]?.isLearning == true && !it.isLearning)) {
                    map[it.id] = it.toJoinedWordProgress()
                }
            }
            val result = map.values.toList()
            Timber.d("getPerWordOverallProgress: $result")
            result
        }
    }
}