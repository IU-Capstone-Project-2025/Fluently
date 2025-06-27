package ru.fluentlyapp.fluently.auth.api

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

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