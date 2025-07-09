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
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

private val ONGOING_LESSON_JSON_KEY = stringPreferencesKey("ongoing_lesson_json_key")

private val lessonModule = SerializersModule {
    polymorphic(LessonComponent::class) {
        subclass(Exercise.InputWord::class)
        subclass(Exercise.FillTheGap::class)
        subclass(Exercise.NewWord::class)
        subclass(Exercise.ChooseTranslation::class)

        subclass(Decoration.Loading::class)
        subclass(Decoration.Finish::class)
        subclass(Decoration.Onboarding::class)
    }
}

private val lessonJsonFormat = Json {
    serializersModule = lessonModule
}

@Singleton
class OngoingLessonDataStore @Inject constructor(
    private val dataStore: DataStore<Preferences>
) {
    /**
     * Get the saved lesson or null if not lesson is saved.
     *
     * May throw exception if something goes wrong with lesson serialization/deserialization.
     */
    fun getOngoingLesson(): Flow<Lesson?> {
        return dataStore.data.map {
            it[ONGOING_LESSON_JSON_KEY]?.let { lessonJson ->
                val decodedLesson: Lesson? = lessonJsonFormat.decodeFromString(lessonJson)
                Timber.d("Decode from saved json: $decodedLesson")
                decodedLesson
            }
        }
    }

    /**
     * Save the lesson.
     *
     * May throw exception if something goes front with lesson serialization/deserialization
     */
    suspend fun setOngoingLesson(lesson: Lesson) {
        dataStore.edit {
            val encodedLesson = lessonJsonFormat.encodeToString(lesson)
            it[ONGOING_LESSON_JSON_KEY] = encodedLesson
            Timber.d("Save ongoing lesson: $encodedLesson")
        }
    }

    /**
     * Clear the ongoing lesson if any is stored.
     */
    suspend fun dropOngoingLesson() {
        dataStore.edit {
            it.remove(ONGOING_LESSON_JSON_KEY)
            Timber.d("Remove ongoing lesson")
        }
    }

}