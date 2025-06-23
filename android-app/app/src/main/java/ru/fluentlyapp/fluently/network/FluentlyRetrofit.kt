package ru.fluentlyapp.fluently.network

import retrofit2.Retrofit
import retrofit2.http.Body
import retrofit2.http.GET
import ru.fluentlyapp.fluently.model.ServerToken
import ru.fluentlyapp.fluently.network.model.RefreshServerTokenRequestBody
import ru.fluentlyapp.fluently.network.model.GetServerTokenRequestBody

private interface FluentlyRetrofitApi {
    @GET("/auth/google")
    fun getServerToken(@Body serverTokenRequestBody: GetServerTokenRequestBody)

    @GET("/auth/refresh")
    fun refreshServerToken(@Body refreshServerTokenRequestBody: RefreshServerTokenRequestBody)
}

private const val FLUENTLY_API_BASE_URL = "example.com" // TODO: change this

class FluentlyNetworkRetrofit : FluentlyNetworkDataSource {
    private val fluentlyRetrofitApi = Retrofit
        .Builder()
        .baseUrl(FLUENTLY_API_BASE_URL)

    override suspend fun getServerToken(idToken: String): ServerToken {
        TODO("Not yet implemented")
    }

    override suspend fun refreshServerToken(refreshToken: String): ServerToken {
        TODO("Not yet implemented")
    }

}

