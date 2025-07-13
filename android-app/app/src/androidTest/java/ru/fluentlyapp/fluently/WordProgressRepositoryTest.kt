package ru.fluentlyapp.fluently

import android.content.Context
import android.util.Log
import androidx.room.Room
import androidx.test.platform.app.InstrumentationRegistry
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Test
import org.junit.Assert.*
import ru.fluentlyapp.fluently.database.app.AppDatabase
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgress
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepositoryImpl
import java.time.OffsetDateTime

class WordProgressRepositoryTest {
    lateinit var wordProgressRepository: WordProgressRepositoryImpl

    @Before
    fun setup() {
        val targetContext = InstrumentationRegistry.getInstrumentation().targetContext
        val appDatabase = Room.inMemoryDatabaseBuilder(
            targetContext,
            AppDatabase::class.java
        ).build()
        wordProgressRepository = WordProgressRepositoryImpl(
            appDatabase.wordProgressDao()
        )
    }

    @Test
    fun basic_crud_works() {
        val wordProgresses = (1..5).map {
            WordProgress(
                it.toString(),
                isLearning = it % 2 == 0,
                instant = OffsetDateTime.now().plusYears(it.toLong()).toInstant()
            )
        }

        runBlocking {
            wordProgresses.forEach { wordProgressRepository.addProgress(it) }
        }

        var expectedProgresses = wordProgresses.slice(1..3).toSet()
        var actualProgresses = runBlocking {
            wordProgressRepository.getProgresses(
                wordProgresses[1].instant,
                wordProgresses[3].instant
            ).first().toSet()
        }
        Log.i("Test", "${expectedProgresses.joinToString()} ${actualProgresses.joinToString()}")
        assertEquals(expectedProgresses, actualProgresses)

        runBlocking {
            wordProgressRepository.removeProgress(wordProgresses[1])
        }

        expectedProgresses = wordProgresses.slice(2..3).toSet()
        actualProgresses = runBlocking {
            wordProgressRepository.getProgresses(
                wordProgresses[1].instant,
                wordProgresses[3].instant
            ).first().toSet()
        }
        Log.i("Test", "${expectedProgresses.joinToString()} ${actualProgresses.joinToString()}")
        assertEquals(expectedProgresses, actualProgresses)
    }
}