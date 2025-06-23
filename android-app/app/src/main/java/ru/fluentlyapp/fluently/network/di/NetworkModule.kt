package ru.fluentlyapp.fluently.network.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.network.FluentlyNetworkDataSource
import ru.fluentlyapp.fluently.network.FluentlyRetrofit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class NetworkModule {
    @Binds
    @Singleton
    abstract fun provideFluentlyNetworkDataSource(
        fluentlyNetworkDataSource: FluentlyNetworkDataSource
    ): FluentlyRetrofit
}