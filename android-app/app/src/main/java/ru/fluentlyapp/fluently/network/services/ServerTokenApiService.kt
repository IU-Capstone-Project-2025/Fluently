package ru.fluentlyapp.fluently.network.services

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST
import ru.fluentlyapp.fluently.network.model.GetServerTokenRequestBody
import ru.fluentlyapp.fluently.network.model.RefreshServerTokenRequest
import ru.fluentlyapp.fluently.network.model.ServerTokenResponseBody

interface ServerTokenApiService {
    @POST("/auth/refresh")
    suspend fun refreshToken(
        @Body refreshServerTokenRequest: RefreshServerTokenRequest
    ): Response<ServerTokenResponseBody>

    @POST("/auth/google")
    suspend fun getServerToken(
        @Body serverTokenRequestBody: GetServerTokenRequestBody
    ): Response<ServerTokenResponseBody>
}
