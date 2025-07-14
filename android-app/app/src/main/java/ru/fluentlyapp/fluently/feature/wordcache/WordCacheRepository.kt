package ru.fluentlyapp.fluently.feature.wordcache

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.database.app.wordcache.WordCacheDao
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

interface WordCacheRepository {
    suspend fun updateWord(wordCache: WordCache)
    suspend fun removeWord(wordCache: WordCache)
    fun getAllWords(): Flow<List<WordCache>>
    fun getWordCacheById(id: String): Flow<WordCache>
}

@Singleton
class WordCacheRepositoryImpl @Inject constructor(
    appDatabase: AppDatabase
): WordCacheRepository {
    private val wordCacheDao = appDatabase.wordCacheDao()
    override suspend fun updateWord(wordCache: WordCache) {
        Timber.d("updateWord: $wordCache")
        wordCacheDao.insert(wordCache.toWordCacheEntity())
    }

    override suspend fun removeWord(wordCache: WordCache) {
        Timber.d("removeWord: $wordCache")
        wordCacheDao.delete(wordCache.toWordCacheEntity())
    }

    override fun getAllWords(): Flow<List<WordCache>> {
        return wordCacheDao.getAll().map { list ->
            val result = list.map { it.toWordCache() }
            Timber.d("getAllWords: $result")
            result
        }
    }

    override fun getWordCacheById(id: String): Flow<WordCache> {
        val result = wordCacheDao.getById(id).map { it.toWordCache() }
        Timber.d("getWordCacheById: id=$id; result=$result")
        return result
    }
}