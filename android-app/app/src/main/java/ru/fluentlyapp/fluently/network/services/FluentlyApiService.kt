package ru.fluentlyapp.fluently.network.services

import retrofit2.Response
import retrofit2.http.GET
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody

interface FluentlyApiService {
    @GET("/api/v1/lesson")
    suspend fun getLesson(): Response<LessonResponseBody>
}