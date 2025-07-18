package ru.fluentlyapp.fluently.feature.joinedwordprogress.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepository
import ru.fluentlyapp.fluently.feature.joinedwordprogress.JoinedWordProgressRepositoryImpl
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class JoinedWordProgressModule {
    @Binds
    @Singleton
    abstract fun bindJoinedWordProgressRepository(
        joinedWordProgressRepositoryImpl: JoinedWordProgressRepositoryImpl
    ): JoinedWordProgressRepository
}