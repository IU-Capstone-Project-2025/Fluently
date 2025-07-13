package ru.fluentlyapp.fluently.datastore

import android.net.Uri
import androidx.core.net.toUri
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject

private val PROFILE_PICTURE_URI = stringPreferencesKey("profile_picture_uri")

class UserPreferencesDataStore @Inject constructor(
  val dataStore: DataStore<Preferences>
) {
    suspend fun setUserProfileUri(uri: Uri) {
        dataStore.edit {
            it[PROFILE_PICTURE_URI] = uri.toString()
        }
    }

    fun getUserProfileUri(): Flow<Uri?> {
        return dataStore.data.map {
            it[PROFILE_PICTURE_URI]?.toUri()
        }
    }

    suspend fun dropUserProfileUri() {
        dataStore.edit {
            it.remove(PROFILE_PICTURE_URI)
        }
    }
}