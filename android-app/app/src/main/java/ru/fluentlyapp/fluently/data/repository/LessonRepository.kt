package ru.fluentlyapp.fluently.data.repository

import ru.fluentlyapp.fluently.model.Lesson

interface LessonRepository {
    suspend fun getCurrentLesson(): Lesson?
    suspend fun setCurrentLesson(lesson: Lesson)
    suspend fun fetchLesson(): Lesson
    suspend fun updateLesson(lesson: Lesson)
}

class StubLessonRepository : LessonRepository {
    override suspend fun getCurrentLesson(): Lesson? {
        return null
    }

    override suspend fun setCurrentLesson(lesson: Lesson) {}

    override suspend fun fetchLesson(): Lesson {
        TODO("Not yet implemented")
    }

    override suspend fun updateLesson(lesson: Lesson) {}
}