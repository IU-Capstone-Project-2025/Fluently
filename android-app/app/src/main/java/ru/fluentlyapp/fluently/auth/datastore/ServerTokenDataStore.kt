package ru.fluentlyapp.fluently.auth.datastore

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.serialization.json.Json
import ru.fluentlyapp.fluently.auth.model.ServerToken
import timber.log.Timber
import javax.inject.Inject
import javax.inject.Singleton

val SERVER_TOKEN_KEY = stringPreferencesKey("SERVER_TOKEN_KEY")

@Singleton
class ServerTokenDataStore @Inject constructor(
    private val preferencesDataStore: DataStore<Preferences>
) {
    suspend fun saveServerToken(serverToken: ServerToken) {
        preferencesDataStore.edit {
            it[SERVER_TOKEN_KEY] = Json.encodeToString(serverToken)
            Timber.d("Save server token: %s", serverToken.toString())
        }
    }

    suspend fun getServerToken(): ServerToken? {
        return preferencesDataStore.data.map {
            it[SERVER_TOKEN_KEY]?.let { serverTokenJson ->
                Json.decodeFromString<ServerToken>(
                    serverTokenJson
                )
            }
        }.first()
    }

    suspend fun deleteServerToken() {
        preferencesDataStore.edit {
            it.remove(SERVER_TOKEN_KEY)
            Timber.d("Delete server token")
        }
    }
}
