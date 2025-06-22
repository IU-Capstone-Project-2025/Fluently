package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import ru.fluentlyapp.fluently.model.Lesson

interface LessonRepository {
    suspend fun getOngoingLessonId(): String
    suspend fun setOngoingLessonId(lessonId: String)
    suspend fun dropOngoingLesson()

    suspend fun fetchCurrentLesson(): Lesson
    suspend fun fetchLesson(lessonId: String): Lesson
    suspend fun saveLesson(lesson: Lesson)
    suspend fun getSavedLesson(lessonId: String): Lesson?
    fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?>
}

class StubLessonRepository : LessonRepository {
    override suspend fun getOngoingLessonId(): String {
        TODO("Not yet implemented")
    }

    override suspend fun setOngoingLessonId(lessonId: String) {
        TODO("Not yet implemented")
    }

    override suspend fun dropOngoingLesson() {
        TODO("Not yet implemented")
    }

    override suspend fun fetchCurrentLesson(): Lesson {
        TODO("Not yet implemented")
    }

    override suspend fun fetchLesson(lessonId: String): Lesson {
        TODO("Not yet implemented")
    }

    override suspend fun saveLesson(lesson: Lesson) {
        TODO("Not yet implemented")
    }

    override suspend fun getSavedLesson(lessonId: String): Lesson? {
        TODO("Not yet implemented")
    }

    override fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?> {
        TODO("Not yet implemented")
    }

}