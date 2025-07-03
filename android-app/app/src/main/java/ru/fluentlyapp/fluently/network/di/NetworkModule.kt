package ru.fluentlyapp.fluently.network.di

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
import ru.fluentlyapp.fluently.common.di.BaseOkHttpClient
import ru.fluentlyapp.fluently.network.FLUENTLY_BASE_URL
import ru.fluentlyapp.fluently.network.services.FluentlyApiService
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import ru.fluentlyapp.fluently.network.FluentlyApiDefaultDataSource
import ru.fluentlyapp.fluently.network.middleware.AccessTokenInterceptor
import ru.fluentlyapp.fluently.network.middleware.AuthAuthenticator
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class NetworkModule {
    companion object {
        @Provides
        @Singleton
        fun provideFluentlyApiService(
            @BaseOkHttpClient baseClient: OkHttpClient,
            accessTokenInterceptor: AccessTokenInterceptor,
            authenticator: AuthAuthenticator
        ): FluentlyApiService {
            val customizedClient = baseClient
                .newBuilder()
                .addInterceptor(accessTokenInterceptor)
                .authenticator(authenticator)
                .build()

            val retrofit = Retrofit
                .Builder()
                .client(customizedClient)
                .addConverterFactory(Json.asConverterFactory("application/json; charset=UTF8".toMediaType()))
                .baseUrl(FLUENTLY_BASE_URL)
                .build()

            return retrofit.create(FluentlyApiService::class.java)
        }
    }

    @Binds
    @Singleton
    abstract fun bindFluentlyApiDataSource(
        fluentlyApiDefaultDataSource: FluentlyApiDefaultDataSource
    ): FluentlyApiDataSource
}