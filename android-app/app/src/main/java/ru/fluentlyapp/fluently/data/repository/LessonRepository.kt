package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.update
import ru.fluentlyapp.fluently.datastore.LessonPreferencesDataStore
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import javax.inject.Inject

interface LessonRepository {
    /**
     * Return the id of locally set ongoing lesson.
     *
     * Returns null if none of the lessons are ongoing.
     */
    fun getSavedOngoingLessonIdAsFlow(): Flow<String?>

    /**
     * Locally, set the `lessonId` as the ongoing lesson.
     */
    suspend fun setSavedOngoingLessonId(lessonId: String)

    /**
     * Locally, drop the ongoing lesson id.
     */
    suspend fun dropSavedOngoingLesson()

    /**
     * Get the saved lesson as `Flow` by the `lessonId`. The flow may emit null if
     * none of the lessons are stored under `lessonId`.
     */
    fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?>

    /**
     * At any moment, any user has the current ongoing lesson. This method fetches
     * the currently assigned lesson for this user.
     *
     * May throw exception.
     */
    suspend fun fetchCurrentLesson(): Lesson

    /**
     * Fetch the lesson by the `lessonId` from the server
     *
     * May throw exception.
     */
    suspend fun fetchLesson(lessonId: String): Lesson

    /**
     * Update the `lesson` locally.
     */
    suspend fun saveLesson(lesson: Lesson)

    /**
     * Get the saved lesson by the `lessonId`.
     *
     * Returns null if no lessons are saved under `lessonId`.
     */
    suspend fun getSavedLesson(lessonId: String): Lesson?

    /**
     * Send the lesson to the server so that it stores it.
     *
     * May throw exception.
     */
    suspend fun sendLesson(lesson: Lesson)
}

class StubLessonRepository @Inject constructor(
    val lessonPreferencesDataStore: LessonPreferencesDataStore,
    val fluentlyApiDataSource: FluentlyApiDataSource
) : LessonRepository {
    val lessons = mutableMapOf<String, MutableStateFlow<Lesson?>>()
    override fun getSavedOngoingLessonIdAsFlow(): Flow<String?> {
        return lessonPreferencesDataStore.getOngoingLessonIdAsFlow()
    }

    override suspend fun setSavedOngoingLessonId(lessonId: String) {
        lessonPreferencesDataStore.setOngoingLessonId(lessonId)
    }

    override suspend fun dropSavedOngoingLesson() {
        lessonPreferencesDataStore.dropOngoingLessonId()
    }

    override suspend fun fetchCurrentLesson(): Lesson {
        return fluentlyApiDataSource.getCurrentLesson()
    }

    override suspend fun fetchLesson(lessonId: String): Lesson {
        return fluentlyApiDataSource.getLesson(lessonId)
    }

    override suspend fun sendLesson(lesson: Lesson) {
        TODO("Not yet implemented")
    }

    override suspend fun saveLesson(lesson: Lesson) {
        val stateFlow = lessons.getOrPut(lesson.lessonId) {
            MutableStateFlow(null)
        }
        stateFlow.update { lesson }
    }

    override suspend fun getSavedLesson(lessonId: String): Lesson? {
        return lessons[lessonId]?.value
    }

    override fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?> {
        val stateFlow = lessons.getOrPut(lessonId) {
            MutableStateFlow(null)
        }

        return stateFlow
    }
}