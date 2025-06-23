package ru.fluentlyapp.fluently.datastore

import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.datastore.model.ServerTokenPreference

fun ServerToken.toServerTokenPreference() = ServerTokenPreference(
    accessToken = accessToken,
    refreshToken = refreshToken,
    expiresInSeconds = expiresInSeconds,
    tokenType = tokenType
)

fun ServerTokenPreference.toServerToken() = ServerToken(
    accessToken = accessToken,
    refreshToken = refreshToken,
    expiresInSeconds = expiresInSeconds,
    tokenType = tokenType
)