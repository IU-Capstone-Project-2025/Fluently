package ru.fluentlyapp.fluently.network.services

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.PUT
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.network.model.internal.CardApiModel
import ru.fluentlyapp.fluently.network.model.internal.ChatRequestBody
import ru.fluentlyapp.fluently.network.model.internal.ChatResponseBody
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.internal.TopicApiModel
import ru.fluentlyapp.fluently.network.model.internal.UserPreferencesRequestBody
import ru.fluentlyapp.fluently.network.model.internal.UserPreferencesResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordOfTheDayResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordProgressApiModel

interface FluentlyApiService {
    @GET("/api/v1/lesson")
    suspend fun getLesson(): Response<LessonResponseBody>

    @POST("/api/v1/progress")
    suspend fun sendProgress(@Body progress: List<WordProgressApiModel>)

    @GET("/api/v1/day-word")
    suspend fun getDayOfTheWord(): Response<WordOfTheDayResponseBody>

    @POST("/api/v1/chat")
    suspend fun sendChat(@Body chat: ChatRequestBody): Response<ChatResponseBody>

    @POST("/api/v1/chat/finish")
    suspend fun sendFinish()

    @GET("/api/v1/preferences")
    suspend fun getUserPreferences(): Response<UserPreferencesResponseBody>

    @PUT("/api/v1/preferences")
    suspend fun putUserPreferences(@Body preferences: UserPreferencesRequestBody)

    @GET("/api/v1/topics")
    suspend fun getTopics(): Response<List<TopicApiModel>>
}