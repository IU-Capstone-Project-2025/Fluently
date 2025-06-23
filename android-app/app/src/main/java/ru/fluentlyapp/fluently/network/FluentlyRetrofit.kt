package ru.fluentlyapp.fluently.network

import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import retrofit2.Retrofit
import retrofit2.converter.kotlinx.serialization.asConverterFactory
import retrofit2.http.Body
import retrofit2.http.GET
import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.network.model.RefreshServerTokenRequestBody
import ru.fluentlyapp.fluently.network.model.GetServerTokenRequestBody
import ru.fluentlyapp.fluently.network.model.ServerTokenResponseBody
import javax.inject.Singleton

private interface FluentlyRetrofitApi {
    @GET("/auth/google")
    suspend fun getServerToken(
        @Body serverTokenRequestBody: GetServerTokenRequestBody
    ): ServerTokenResponseBody

    @GET("/auth/refresh")
    fun refreshServerToken(
        @Body refreshServerTokenRequestBody: RefreshServerTokenRequestBody
    ): ServerTokenResponseBody
}

private const val FLUENTLY_API_BASE_URL = "example.com"

@Singleton
class FluentlyRetrofit : FluentlyNetworkDataSource {
    private val retrofit = Retrofit
        .Builder()
        .addConverterFactory(
            Json.asConverterFactory(
                "application/json; charset=UTF8".toMediaType(),
            )
        )
        .baseUrl(FLUENTLY_API_BASE_URL)
        .build()

    private val fluentlyRetrofitApi = retrofit.create(FluentlyRetrofitApi::class.java)

    override suspend fun getServerToken(idToken: String): ServerToken {
        val result = fluentlyRetrofitApi.getServerToken(
            GetServerTokenRequestBody(
                idToken = idToken,
                platform = "android"
            )
        )
        return result.toServerToken()
    }

    override suspend fun refreshServerToken(refreshToken: String): ServerToken {
        val result = fluentlyRetrofitApi.refreshServerToken(
            refreshServerTokenRequestBody = RefreshServerTokenRequestBody(
                refreshToken = refreshToken
            )
        )
        return result.toServerToken()
    }

}

