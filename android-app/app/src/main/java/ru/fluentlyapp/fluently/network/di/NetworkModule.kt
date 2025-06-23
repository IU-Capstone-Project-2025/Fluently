package ru.fluentlyapp.fluently.network.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.serialization.json.Json
import okhttp3.MediaType
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.kotlinx.serialization.asConverterFactory
import ru.fluentlyapp.fluently.network.FLUENTLY_BASE_URL
import ru.fluentlyapp.fluently.network.services.ServerTokenApiService
import javax.inject.Qualifier
import javax.inject.Singleton

@Qualifier
@Retention(AnnotationRetention.RUNTIME)
annotation class BaseClient


@Module
@InstallIn(SingletonComponent::class)
class NetworkModule {
    @Provides
    @Singleton
    @BaseClient
    fun provideBaseClient(): OkHttpClient {
        val loggingInterceptor = HttpLoggingInterceptor()
        loggingInterceptor.level = HttpLoggingInterceptor.Level.BODY
        return OkHttpClient
            .Builder()
            .addInterceptor(loggingInterceptor)
            .build()
    }

    @Provides
    @Singleton
    fun provideServerTokenService(
        @BaseClient baseClient: OkHttpClient
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