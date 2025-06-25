package ru.fluentlyapp.fluently.datastore

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

private val ONGOING_LESSON_ID_KEY = stringPreferencesKey("ongoing_lesson_id_key")

@Singleton
class LessonPreferencesDataStore @Inject constructor(
    private val dataStore: DataStore<Preferences>
) {
    suspend fun setOngoingLessonId(lessonId: String) {
        dataStore.edit {
            it[ONGOING_LESSON_ID_KEY] = lessonId
        }
    }

    suspend fun getOngoingLessonId(): String? = dataStore.data.map {
        it[ONGOING_LESSON_ID_KEY]
    }.first()

    suspend fun dropOngoingLessonId() {
        dataStore.edit {
            it.remove(ONGOING_LESSON_ID_KEY)
        }
    }
}