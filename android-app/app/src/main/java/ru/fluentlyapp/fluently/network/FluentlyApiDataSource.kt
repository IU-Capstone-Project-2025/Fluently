package ru.fluentlyapp.fluently.network

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.supervisorScope
import kotlinx.coroutines.withContext
import retrofit2.HttpException
import retrofit2.Response
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.network.model.Chat
import ru.fluentlyapp.fluently.network.model.Progress
import ru.fluentlyapp.fluently.network.model.WordOfTheDay
import ru.fluentlyapp.fluently.network.model.internal.CardApiModel
import ru.fluentlyapp.fluently.network.model.internal.ChatResponseBody
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordOfTheDayResponseBody
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
import timber.log.Timber
import javax.inject.Inject

interface FluentlyApiDataSource {
    suspend fun getLesson(): Lesson

    suspend fun sendProgress(progress: Progress): Unit

    suspend fun getWordOfTheDay(): WordOfTheDay

    suspend fun sendChat(chat: Chat): Chat

    suspend fun sendFinish()
}

class FluentlyApiDefaultDataSource @Inject constructor(
    private val fluentlyApiService: FluentlyApiService
) : FluentlyApiDataSource {
    private fun <T> getSuccessfulResponseBody(response: Response<T>): T {
        if (!response.isSuccessful) {
            throw HttpException(response)
        }

        val responseBody = response.body()

        if (responseBody == null) {
            throw IllegalStateException("responseBody is null")
        }

        return responseBody
    }

    override suspend fun getLesson(): Lesson {
        return withContext(Dispatchers.IO) {
            Timber.d("Performing request for lesson from `fluentlyApiService`")
            val response = fluentlyApiService.getLesson()
            Timber.d("Receive response from the fluentlyApiService; code=${response.code()}; message=${response.message()}")
            val body: LessonResponseBody = getSuccessfulResponseBody(response)
            body.convertToLesson()
        }
    }

    override suspend fun sendProgress(progress: Progress) {
        withContext(Dispatchers.IO) {
            Timber.d("Performing request sendProgress: $progress")

            fluentlyApiService.sendProgress(progress.toProgressRequestBody())
        }
    }

    override suspend fun getWordOfTheDay(): WordOfTheDay {
        return withContext(Dispatchers.IO) {
            Timber.d("Performing request getWordOfTheDay")
            val response = fluentlyApiService.getDayOfTheWord()
            val body: WordOfTheDayResponseBody = getSuccessfulResponseBody(response)
            body.toWordOfTheDay()
        }
    }

    override suspend fun sendChat(chat: Chat): Chat {
        return withContext(Dispatchers.IO) {
            Timber.d("Performing request sendChat")
            val response = fluentlyApiService.sendChat(chat.toChatRequestBody())
            val body: ChatResponseBody = getSuccessfulResponseBody(response)
            body.toChat()
        }
    }

    override suspend fun sendFinish() {
        withContext(Dispatchers.IO) {
            Timber.d("Performing sendFinish")
            fluentlyApiService.sendFinish()
        }
    }
}
