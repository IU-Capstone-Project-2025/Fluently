package ru.fluentlyapp.fluently.feature.userpreferences

import kotlinx.coroutines.flow.Flow
import ru.fluentlyapp.fluently.common.model.UserPreferences
import ru.fluentlyapp.fluently.datastore.UserPreferencesDataStore
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class UserPreferencesRepository @Inject constructor(
    private val fluentlyApiDataSource: FluentlyApiDataSource,
    private val userPreferencesDataStore: UserPreferencesDataStore
) {
    suspend fun updateUserPreferences() {
        val userPreferences = fluentlyApiDataSource.getUserPreferences()
        userPreferencesDataStore.setUserPreferences(userPreferences)
        Timber.v("Update user preferences: $userPreferences")
    }

    fun getUserPreferences(): Flow<UserPreferences?> {
        return userPreferencesDataStore.getUserPreferences()
    }

    suspend fun dropUserPreferences() {
        userPreferencesDataStore.dropUserPreferences()
    }
}