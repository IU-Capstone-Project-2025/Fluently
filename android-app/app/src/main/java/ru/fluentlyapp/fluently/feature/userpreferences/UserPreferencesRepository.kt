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
    suspend fun updateCachedUserPreferences(userPreferences: UserPreferences) {
        Timber.v("Update user preferences: $userPreferences")
        userPreferencesDataStore.setUserPreferences(userPreferences)
    }

    suspend fun getRemoteUserPreferences(): UserPreferences {
        val result = fluentlyApiDataSource.getUserPreferences()
        Timber.d("getRemoteUserPreferences: $result")
        return result
    }

    fun getCachedUserPreferences(): Flow<UserPreferences?> {
        return userPreferencesDataStore.getUserPreferences()
    }

    suspend fun updateRemoteUserPreferences(preferences: UserPreferences) {
        Timber.d("updateRemoveUserPreferences: $preferences")
        fluentlyApiDataSource.sendUserPreferences(preferences)
    }

    suspend fun dropCachedUserPreferences() {
        Timber.d("dropCachedUserPreferences")
        userPreferencesDataStore.dropUserPreferences()
    }
}