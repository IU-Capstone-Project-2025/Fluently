package ru.fluentlyapp.fluently.auth.api

import ru.fluentlyapp.fluently.auth.model.ServerToken

fun ServerTokenResponseBody.toServerToken() = ServerToken(
    accessToken = accessToken,
    refreshToken = refreshToken,
    tokenType = tokenType,
    expiresInSeconds = expiresInSeconds
)