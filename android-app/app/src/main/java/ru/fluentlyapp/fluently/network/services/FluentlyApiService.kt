package ru.fluentlyapp.fluently.network.services

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordProgressApiModel

interface FluentlyApiService {
    @GET("/api/v1/lesson")
    suspend fun getLesson(): Response<LessonResponseBody>

    @POST("/api/v1/progress")
    suspend fun sendProgress(@Body progress: List<WordProgressApiModel>)
}