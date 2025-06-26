package ru.fluentlyapp.fluently

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.emptyPreferences
import androidx.datastore.preferences.preferencesDataStore
import androidx.test.ext.junit.runners.AndroidJUnit4
import androidx.test.platform.app.InstrumentationRegistry
import kotlinx.coroutines.runBlocking
import org.junit.After
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith
import ru.fluentlyapp.fluently.auth.model.ServerToken
import ru.fluentlyapp.fluently.auth.datastore.ServerTokenDataStore

// Code here was partially generated using AI tools
@RunWith(AndroidJUnit4::class)
class ServerTokenDataStoreTest {
    private val Context.testDataStore by preferencesDataStore(
        name = "test_preferences_${System.currentTimeMillis()}"
    )

    private lateinit var context: Context
    private lateinit var serverTokenDataStore: ServerTokenDataStore

    @Before
    fun setUp() {
        context = InstrumentationRegistry.getInstrumentation().targetContext
        serverTokenDataStore = ServerTokenDataStore(context.testDataStore)
    }

    suspend fun clearDataStore() {
        context.testDataStore.edit { emptyPreferences() }
    }

    @After
    fun cleanUp() = runBlocking {
        clearDataStore()
    }

    @Test
    fun saveAndRetrieveServerToken_returnsCorrectValue() = runBlocking {
        clearDataStore()

        val serverToken = ServerToken(
            accessToken = "abc123",
            refreshToken = "refresh456",
            tokenType = "Bearer",
            expiresInSeconds = 100
        )

        serverTokenDataStore.saveServerToken(serverToken)
        val retrievedToken = serverTokenDataStore.getServerToken()

        assertEquals(serverToken, retrievedToken)
    }

    @Test
    fun getServerToken_whenNoTokenSaved_returnsNull() = runBlocking {
        clearDataStore()
        val token = serverTokenDataStore.getServerToken()
        assertNull(token)
    }

    @Test
    fun whenServerTokenIsDeleted_tokenIsRemovedFromDataStore() = runBlocking {
        clearDataStore()

        val serverToken = ServerToken(
            accessToken = "abc123",
            refreshToken = "refresh456",
            tokenType = "Bearer",
            expiresInSeconds = 100
        )

        serverTokenDataStore.saveServerToken(serverToken)
        serverTokenDataStore.deleteServerToken()

        val tokenAfterDelete = serverTokenDataStore.getServerToken()
        assertNull(tokenAfterDelete)
    }
}