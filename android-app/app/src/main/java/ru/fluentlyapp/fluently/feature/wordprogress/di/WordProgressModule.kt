package ru.fluentlyapp.fluently.feature.wordprogress.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepository
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepositoryImpl
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class WordProgressModule {
    @Binds
    @Singleton
    abstract fun bindWordProgressRepository(
        wordProgressRepositoryImpl: WordProgressRepositoryImpl
    ): WordProgressRepository
}