package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.network.model.ServerTokenResponseBody

fun ServerTokenResponseBody.toServerToken() = ServerToken(
    accessToken = accessToken,
    refreshToken = refreshToken,
    tokenType = tokenType,
    expiresInSeconds = expiresInSeconds
)