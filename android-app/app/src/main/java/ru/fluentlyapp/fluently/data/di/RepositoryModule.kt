package ru.fluentlyapp.fluently.data.di

import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.data.repository.DefaultLessonRepository
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {
    @Binds
    @Singleton
    abstract fun bindLessonRepository(
        defaultLessonRepository: DefaultLessonRepository
    ): LessonRepository
}