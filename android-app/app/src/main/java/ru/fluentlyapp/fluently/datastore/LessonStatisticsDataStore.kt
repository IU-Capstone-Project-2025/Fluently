package ru.fluentlyapp.fluently.datastore

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.common.model.Lesson
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

@Serializable
data class LessonsStatistic(
    val knownWords: Int,
    val wordsInProgress: Int
)

private val KEY_LESSONS_STATISTIC = stringPreferencesKey("key_lessons_statistic")

@Singleton
class LessonsStatisticDataStore @Inject constructor(
    val dataStore: DataStore<Preferences>
) {
    suspend fun setLessonsStatistic(lessonsStatistic: LessonsStatistic) {
        dataStore.edit {
            it[KEY_LESSONS_STATISTIC] = Json.encodeToString(lessonsStatistic)
            Timber.d("Save lessonsStatistics: $lessonsStatistic")
        }
    }

    fun getLessonsStatistic(): Flow<LessonsStatistic?> {
        return dataStore.data.map {
            it[KEY_LESSONS_STATISTIC]?.let { it ->
                val decodedLessonsStatistic: LessonsStatistic? = Json.decodeFromString(it)
                Timber.d("Decode from saved json: $decodedLessonsStatistic")
                decodedLessonsStatistic
            }
        }
    }

    suspend fun dropLessonsStatistic() {
        dataStore.edit {
            it.remove(KEY_LESSONS_STATISTIC)
        }
    }
}