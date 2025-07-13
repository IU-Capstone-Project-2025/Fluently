package ru.fluentlyapp.fluently.feature.wordcache

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.database.app.wordcache.WordCacheDao
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
    private val wordCacheDao: WordCacheDao
): WordCacheRepository {
    override suspend fun updateWord(wordCache: WordCache) {
        wordCacheDao.insert(wordCache.toWordCacheEntity())
    }

    override suspend fun removeWord(wordCache: WordCache) {
        wordCacheDao.delete(wordCache.toWordCacheEntity())
    }

    override fun getAllWords(): Flow<List<WordCache>> {
        return wordCacheDao.getAll().map { list ->
            list.map { it.toWordCache() }
        }
    }

    override fun getWordCacheById(id: String): Flow<WordCache> {
        return wordCacheDao.getById(id).map { it.toWordCache() }
    }
}