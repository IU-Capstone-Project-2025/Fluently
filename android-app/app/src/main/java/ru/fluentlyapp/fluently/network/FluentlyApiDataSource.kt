package ru.fluentlyapp.fluently.network

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import retrofit2.HttpException
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
import timber.log.Timber
import javax.inject.Inject

interface FluentlyApiDataSource {
    /**
     * Fetch the generated lesson user from the server.
     *
     * May throw exception.
     */
    suspend fun getLesson(): Lesson
}

class FluentlyApiDefaultDataSource @Inject constructor(
    private val fluentlyApiService: FluentlyApiService
) : FluentlyApiDataSource {
    override suspend fun getLesson(): Lesson {
        return withContext(Dispatchers.IO) {
            Timber.d("Performing request for lesson from `fluentlyApiService`")
            val response = fluentlyApiService.getLesson()
            Timber.d("Receive response from the fluentlyApiService; code=${response.code()}; message=${response.message()}")

            if (!response.isSuccessful) {
                throw HttpException(response)
            }

            val responseBody = response.body()

            if (responseBody == null) {
                throw IllegalStateException("responseBody is null")
            }

            responseBody.convertToLesson()
        }
    }
}
