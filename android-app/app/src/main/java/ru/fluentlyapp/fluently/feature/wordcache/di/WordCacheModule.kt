package ru.fluentlyapp.fluently.feature.wordcache.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepository
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepositoryImpl
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class WordCacheModule {
    @Binds
    @Singleton
    abstract fun bindWordCacheRepository(
        wordCacheRepositoryImpl: WordCacheRepositoryImpl
    ): WordCacheRepository
}