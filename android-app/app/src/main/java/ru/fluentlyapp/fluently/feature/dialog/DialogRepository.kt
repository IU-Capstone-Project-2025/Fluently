package ru.fluentlyapp.fluently.feature.dialog

import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import ru.fluentlyapp.fluently.network.model.Chat
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class DialogRepository @Inject constructor(
    private val fluentlyApiDataSource: FluentlyApiDataSource
) {
    suspend fun sendFinish() {
        Timber.d("sendFinish")
        fluentlyApiDataSource.sendFinish()
    }

    suspend fun sendChat(chat: Chat): Chat {
        Timber.d("sendChat")
        val result = fluentlyApiDataSource.sendChat(chat)
        Timber.d("sendChat: $result")
        return result
    }
}