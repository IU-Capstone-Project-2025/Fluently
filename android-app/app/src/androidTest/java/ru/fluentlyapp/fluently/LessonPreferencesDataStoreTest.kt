package ru.fluentlyapp.fluently

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.emptyPreferences
import androidx.datastore.preferences.preferencesDataStore
import androidx.test.platform.app.InstrumentationRegistry
import kotlinx.coroutines.runBlocking
import org.junit.After
import org.junit.Before
import org.junit.Assert.*
import org.junit.Test
import ru.fluentlyapp.fluently.datastore.LessonPreferencesDataStore


class LessonPreferencesDataStoreTest {
    private val Context.testDataStore by preferencesDataStore(
        name = "test_preferences_${System.currentTimeMillis()}"
    )

    private lateinit var context: Context
    private lateinit var lessonPreferencesDataStore: LessonPreferencesDataStore

    @Before
    fun setUp() {
        context = InstrumentationRegistry.getInstrumentation().targetContext
        lessonPreferencesDataStore = LessonPreferencesDataStore(context.testDataStore)
    }

    suspend fun clearDataStore() {
        context.testDataStore.edit { emptyPreferences() }
    }

    @After
    fun cleanUp() = runBlocking {
        clearDataStore()
    }


    @Test
    fun getLessonIdAfterSetLessonId_returnsCorrectLessonId() = runBlocking {
        clearDataStore()
        val lessonId = "lesson_123"

        lessonPreferencesDataStore.setOngoingLessonId(lessonId)
        val retrievedLessonId = lessonPreferencesDataStore.getOngoingLessonId()

        assertEquals(lessonId, retrievedLessonId)
    }

    @Test
    fun getLessonIdAfterDropLessonId_returnsNull() = runBlocking {
        clearDataStore()
        val lessonId = "lesson_456"

        lessonPreferencesDataStore.setOngoingLessonId(lessonId)
        lessonPreferencesDataStore.dropOngoingLessonId()
        val retrievedLessonId = lessonPreferencesDataStore.getOngoingLessonId()

        assertNull(retrievedLessonId)
    }
}