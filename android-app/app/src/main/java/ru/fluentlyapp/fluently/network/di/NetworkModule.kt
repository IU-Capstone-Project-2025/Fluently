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
import retrofit2.Response
import retrofit2.Retrofit
import retrofit2.converter.kotlinx.serialization.asConverterFactory
import ru.fluentlyapp.fluently.network.FLUENTLY_BASE_URL
import ru.fluentlyapp.fluently.network.middleware.AccessTokenInterceptor
import ru.fluentlyapp.fluently.network.middleware.AuthAuthenticator
import ru.fluentlyapp.fluently.network.model.LessonResponseBody
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
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

//    @Provides
//    @Singleton
//    fun provideFluentlyApiService(
//        @BaseClient baseClient: OkHttpClient,
//        accessTokenInterceptor: AccessTokenInterceptor,
//        authenticator: AuthAuthenticator
//    ): FluentlyApiService {
//        val customizedClient = baseClient
//            .newBuilder()
//            .addInterceptor(accessTokenInterceptor)
//            .authenticator(authenticator)
//            .build()
//
//        val retrofit = Retrofit
//            .Builder()
//            .client(customizedClient)
//            .addConverterFactory(Json.asConverterFactory("application/json; charset=UTF8".toMediaType()))
//            .baseUrl(FLUENTLY_BASE_URL)
//            .build()
//
//        return retrofit.create(FluentlyApiService::class.java)
//    }

    @Provides
    @Singleton
    fun provideMockFluentlyApiService(): FluentlyApiService {
        return object : FluentlyApiService {
            override suspend fun getLesson(): Response<LessonResponseBody> {
                return Response.success<LessonResponseBody>(mockLessonResponse)
            }
        }
    }
}