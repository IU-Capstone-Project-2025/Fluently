package ru.fluentlyapp.fluently.network

import kotlinx.coroutines.delay
import retrofit2.HttpException
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
import ru.fluentlyapp.fluently.testing.mockLessonResponse
import javax.inject.Inject
import javax.inject.Singleton
import kotlin.time.Duration.Companion.milliseconds

@Singleton
class FluentlyApiDataSource @Inject constructor(
    private val fluentlyApiService: FluentlyApiService
) {
    suspend fun getLesson(lessonId: String): Lesson {
        if (lessonId == mockLessonResponse.lesson.lesson_id) {
            return mockLessonResponse.convertToLesson()
        }

        throw IllegalStateException(
            "Invalid lessonId; passed: $lessonId; expected: ${mockLessonResponse.lesson.lesson_id}"
        )
    }

    suspend fun getCurrentLesson(): Lesson {
        delay(2000.milliseconds) // Simulate delay

        val response = fluentlyApiService.getLesson()

        if (!response.isSuccessful) {
            throw HttpException(response)
        }

        val responseBody = response.body()

        if (responseBody == null) {
            throw IllegalStateException("responseBody is null")
        }

        return responseBody.convertToLesson()
    }
}