package ru.fluentlyapp.fluently.network

import retrofit2.HttpException
import ru.fluentlyapp.fluently.model.Lesson
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class FluentlyDataSource @Inject constructor(
    private val fluentlyApiService: FluentlyApiService
) {
    suspend fun getCurrentLesson(): Lesson {
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