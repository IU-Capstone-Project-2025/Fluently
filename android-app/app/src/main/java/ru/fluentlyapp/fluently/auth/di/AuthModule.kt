package ru.fluentlyapp.fluently.auth.di

import dagger.Binds
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.kotlinx.serialization.asConverterFactory
import ru.fluentlyapp.fluently.auth.AuthManager
import ru.fluentlyapp.fluently.auth.GoogleBasedOAuthManager
import ru.fluentlyapp.fluently.auth.api.ServerTokenApiService
import ru.fluentlyapp.fluently.common.di.BaseOkHttpClient
import javax.inject.Singleton

const val FLUENTLY_BASE_URL = "https://fluently-app.ru"

@Module
@InstallIn(SingletonComponent::class)
abstract class AuthModule {
    companion object {
        @Provides
        @Singleton
        fun provideServerTokenApiService(
            @BaseOkHttpClient baseClient: OkHttpClient
        ): ServerTokenApiService {
            val retrofit = Retrofit
                .Builder()
                .client(baseClient)
                .addConverterFactory(Json.asConverterFactory("application/json; charset=UTF8".toMediaType()))
                .baseUrl(FLUENTLY_BASE_URL)
                .build()

            return retrofit.create(ServerTokenApiService::class.java)
        }
    }


    @Binds
    @Singleton
    abstract fun bindAuthRepository(
        googleBasedAuthRepository: GoogleBasedOAuthManager
    ): AuthManager
}