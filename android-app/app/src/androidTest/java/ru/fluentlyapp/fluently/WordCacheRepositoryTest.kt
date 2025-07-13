package ru.fluentlyapp.fluently

import androidx.room.Room
import androidx.test.platform.app.InstrumentationRegistry
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Test
import org.junit.Assert.*
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepository
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepositoryImpl

class WordCacheRepositoryTest {
    lateinit var wordCacheRepository: WordCacheRepository

    @Before
    fun setup() {
        val targetContext = InstrumentationRegistry.getInstrumentation().targetContext
        val appDatabase = Room.inMemoryDatabaseBuilder(
            targetContext,
            AppDatabase::class.java
        ).build()
        wordCacheRepository = WordCacheRepositoryImpl(appDatabase.wordCacheDao())
    }

    @Test
    fun basic_crud_works() {
        val words = (1..5).map {
            WordCache(
                wordId = it.toString(),
                word = "word #$it",
                translation = "перевод #$it",
                examples = listOf("$it a" to "$it b", "$it c" to "$it d")
            )
        }
        runBlocking {
            words.forEach {
                wordCacheRepository.updateWord(it)
            }
        }
        var expected = words.toSet()
        var actual = runBlocking {
            wordCacheRepository.getAllWords().first().toSet()
        }
        assertEquals(expected, actual)

        var expectedWord = words[0]
        var actualWord = runBlocking {
            wordCacheRepository.getWordCacheById(words[0].wordId).first()
        }
        assertEquals(expectedWord, actualWord)

        runBlocking {
            wordCacheRepository.removeWord(words[0])
        }

        expected = words.slice(1..words.size - 1).toSet()
        actual = runBlocking {
            wordCacheRepository.getAllWords().first().toSet()
        }
        assertEquals(expected, actual)
    }
}