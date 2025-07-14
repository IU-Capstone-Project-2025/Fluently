package ru.fluentlyapp.fluently.feature.wordprogress

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.database.app.wordprogress.WordProgressDao
import timber.log.Timber
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

interface WordProgressRepository {
    suspend fun addProgress(wordProgress: WordProgress)
    fun getProgresses(beginDate: Instant, endDate: Instant): Flow<List<WordProgress>>
}

@Singleton
class WordProgressRepositoryImpl @Inject constructor(
    appDatabase: AppDatabase
) : WordProgressRepository {
    private val wordProgressDao = appDatabase.wordProgressDao()

    override suspend fun addProgress(wordProgress: WordProgress) {
        Timber.d("addProgress $wordProgress")
        wordProgressDao.insert(wordProgress.toWordProgressEntity())
    }

    override fun getProgresses(
        beginDate: Instant,
        endDate: Instant
    ): Flow<List<WordProgress>> {
        val result = wordProgressDao.getProgressesBetweenDates(
            beginDate,
            endDate
        ).map {
            it.map { it.toWordProgress() }
        }
        Timber.v("getProgresses: $result")
        return result
    }
}