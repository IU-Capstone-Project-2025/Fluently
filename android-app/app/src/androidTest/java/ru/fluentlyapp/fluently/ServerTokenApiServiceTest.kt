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
    accessToken = "abacaba",
    refreshToken = "opa gangnam style",
    idToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg4MjUwM2E1ZmQ1NmU5ZjczNGRmYmE1YzUwZDdiZjQ4ZGIyODRhZTkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI1MDU5MDM1MTI5ODgtc2hrMzQyaGM3M29tY2E0ZG5uZHYzamM5OHIzOGxscXEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI1MDU5MDM1MTI5ODgtc2hrMzQyaGM3M29tY2E0ZG5uZHYzamM5OHIzOGxscXEuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTU1NDg3MDk0MTg5ODE5NzYzMTkiLCJlbWFpbCI6Im5laWx6dmVzdEBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6ImRDeUE0aVZsRmlLWkNsYXNYdDhVZ1EiLCJub25jZSI6InpzNUJucUtySE1EaG0waG1tdGpUT3ciLCJpYXQiOjE3NTA3MTk2NjAsImV4cCI6MTc1MDcyMzI2MH0.MSUshpEl7JJG_bLnHcVuVirtj_ZpWg7RQmK4PcXhUVsDLFBR3n0k3mczVEnzr_UiV2Sdw9WdOtCA8U1CtUjqfe0jEQMeFUyxY-Qp6x-ozuYNchLtmoxbOEAYQ8wrLcNuViZbmE3uuZkfKFT5lRoQVH8wRZBDevctLZby-tZGjhsZQX1f_slDyn040h0mJE2qg7V-ZkoSAwThipuTnCq3XBbrdvqUuMLp35JfHlgCND9rLxEL-BUNkRhPpry3d6lyL_Bna9K8k_SZVRlSks0B7kfbtprosbYvuvsqHWNkw_dhKZV_Qp4-G7SeoFH4xoZox0l9_ty_ZrpEbVt2zHEnzg"
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

        Log.i("ServerTokenApiServiceTest", "${response.isSuccessful}, ${response.body()}")

        assert(response.isSuccessful && response.body() != null)
    }
}