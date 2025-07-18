package ru.fluentlyapp.fluently.datastore

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import kotlinx.serialization.modules.SerializersModule
import kotlinx.serialization.modules.polymorphic
import kotlinx.serialization.modules.subclass
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.network.model.WordOfTheDay
import timber.log.Timber
import java.time.Instant
import javax.inject.Inject
import javax.inject.Singleton

private val WORD_OF_THE_DAY = stringPreferencesKey("word_of_the_day")
private val LAST_SAVED = stringPreferencesKey("last_saved")

@Singleton
class WordOfTheDayDatastore @Inject constructor(
    private val dataStore: DataStore<Preferences>
) {
    fun getSavedWordOfTheDay(): Flow<WordOfTheDay?> {
        return dataStore.data.map {
            it[WORD_OF_THE_DAY]?.let { wordOfTheDayJson ->
                val decoded: WordOfTheDay = Json.decodeFromString(wordOfTheDayJson)
                Timber.d("Decode word of the day: $decoded")
                decoded
            }
        }
    }

    suspend fun setWordOfTheDay(wordOfTheDay: WordOfTheDay) {
        dataStore.edit {
            val encoded = Json.encodeToString(wordOfTheDay)
            it[WORD_OF_THE_DAY] = encoded
            Timber.d("Save word of the day: $encoded")
            it[LAST_SAVED] = Instant.now().toString()
        }
    }

    fun getLastSaved(): Flow<Instant?> {
        return dataStore.data.map {
            it[LAST_SAVED]?.let { lastSaved ->
                val decoded: Instant = Instant.parse(lastSaved)
                Timber.d("Decode last saved: $decoded")
                decoded
            }
        }
    }

    suspend fun dropWordOfTheDay() {
        dataStore.edit {
            it.remove(WORD_OF_THE_DAY)
            it.remove(LAST_SAVED)
            Timber.d("Remove word of the day")
        }
    }

}