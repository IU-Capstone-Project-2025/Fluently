package ru.fluentlyapp.fluently.network.services

import retrofit2.Response
import retrofit2.http.GET
import ru.fluentlyapp.fluently.network.model.LessonResponseBody

interface FluentlyApiService {
    @GET("/lesson")
    suspend fun getLesson(): Response<LessonResponseBody>
}