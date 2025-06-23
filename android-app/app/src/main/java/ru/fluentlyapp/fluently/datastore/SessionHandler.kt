package ru.fluentlyapp.fluently.datastore

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.intPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.datastore.model.ServerTokenPreference
import javax.inject.Inject
import javax.inject.Singleton

val SERVER_TOKEN_KEY = stringPreferencesKey("SERVER_TOKEN_KEY")

@Singleton
class SessionHandler @Inject constructor(
    private val preferencesDataStore: DataStore<Preferences>
) {
    suspend fun saveServerToken(serverToken: ServerToken) {
        preferencesDataStore.edit {
            it[SERVER_TOKEN_KEY] = Json.encodeToString(serverToken.toServerTokenPreference())
        }
    }

    suspend fun getServerToken(): ServerToken {
        return preferencesDataStore.data.map {
            Json.decodeFromString<ServerTokenPreference>(
                it[SERVER_TOKEN_KEY]!!
            ).toServerToken()
        }.first()
    }

    suspend fun deleteServerToken() {
        preferencesDataStore.edit {
            it.remove(SERVER_TOKEN_KEY)
        }
    }
}
