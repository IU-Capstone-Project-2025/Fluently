package ru.fluentlyapp.fluently.datastore

import android.net.Uri
import androidx.core.net.toUri
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.common.model.UserPreferences
import javax.inject.Inject

private val USER_PREFERENCES_KEY = stringPreferencesKey("user_preferences_key")

class UserPreferencesDataStore @Inject constructor(
  val dataStore: DataStore<Preferences>
) {
    suspend fun setUserPreferences(userPreferences: UserPreferences) {
        dataStore.edit {
            it[USER_PREFERENCES_KEY] = Json.encodeToString(userPreferences)
        }
    }

    fun getUserPreferences(): Flow<UserPreferences?> {
        return dataStore.data.map {
            val decoded = it[USER_PREFERENCES_KEY]
            if (decoded == null) {
                return@map null
            }
            return@map Json.decodeFromString<UserPreferences>(decoded)
        }
    }

    suspend fun dropUserPreferences() {
        dataStore.edit {
            it.remove(USER_PREFERENCES_KEY)
        }
    }
}