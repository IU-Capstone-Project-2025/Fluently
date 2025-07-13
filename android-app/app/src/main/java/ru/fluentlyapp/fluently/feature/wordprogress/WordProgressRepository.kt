package ru.fluentlyapp.fluently.feature.wordprogress

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.feature.wordprogress.database.WordProgressDatabase
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

interface WordProgressRepository {
    suspend fun addProgress(wordProgress: WordProgress)
    suspend fun removeProgress(wordProgress: WordProgress)
    fun getProgresses(beginDate: Instant, endDate: Instant): Flow<List<WordProgress>>
}

@Singleton
class WordProgressRepositoryImpl @Inject constructor(
    wordProgressDatabase: WordProgressDatabase
) : WordProgressRepository {
    private val wordProgressDao = wordProgressDatabase.wordProgressDao()

    override suspend fun addProgress(wordProgress: WordProgress) {
        wordProgressDao.insert(wordProgress.toWordProgressEntity())
    }

    override suspend fun removeProgress(wordProgress: WordProgress) {
        wordProgressDao.delete(wordProgress.toWordProgressEntity())
    }

    override fun getProgresses(
        beginDate: Instant,
        endDate: Instant
    ): Flow<List<WordProgress>> {
        return wordProgressDao.getProgressesBetweenDates(
            beginDate,
            endDate
        ).map {
            it.map { it.toWordProgress() }
        }
    }
}