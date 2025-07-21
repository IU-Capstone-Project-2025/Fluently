package ru.fluentlyapp.fluently.feature.topics

import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TopicRepository @Inject constructor(
    private val fluentlyApiDataSource: FluentlyApiDataSource
) {
    suspend fun getAvailableTopic(): List<String> {
        return fluentlyApiDataSource.getTopics()
    }
}