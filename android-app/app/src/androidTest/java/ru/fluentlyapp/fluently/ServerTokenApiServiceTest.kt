package ru.fluentlyapp.fluently

import android.util.Log
import dagger.hilt.android.testing.HiltAndroidRule
import dagger.hilt.android.testing.HiltAndroidTest
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.Assert.*
import ru.fluentlyapp.fluently.network.model.GetServerTokenRequestBody
import ru.fluentlyapp.fluently.network.services.ServerTokenApiService
import ru.fluentlyapp.fluently.oauth.model.OAuthToken
import javax.inject.Inject

// The only hope is idToken can be persisted across tests
val TestOAuthToken = OAuthToken(
    accessToken = "fake access token",
    refreshToken = "fake refresh token",
    idToken = "fake id token"
)

@HiltAndroidTest
class ServerTokenApiServiceTest {
    @get:Rule
    val hiltRule = HiltAndroidRule(this)

    @Inject
    lateinit var serverTokenApiService: ServerTokenApiService

    @Before
    fun setup() {
        hiltRule.inject()
    }

    @Test
    fun getServerToken_returnsServerToken(): Unit = runBlocking {
        val response = serverTokenApiService.getServerToken(
            serverTokenRequestBody = GetServerTokenRequestBody(
                idToken = TestOAuthToken.idToken,
                platform = "android"
            )
        )

        assert(response.code() == 401)
    }
}