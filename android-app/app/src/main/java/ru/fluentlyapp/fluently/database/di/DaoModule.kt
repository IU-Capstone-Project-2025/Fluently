package ru.fluentlyapp.fluently.database.di

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.database.FluentlyDatabase
import ru.fluentlyapp.fluently.database.dao.LessonDao
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
class DaoModule {
    @Provides
    @Singleton
    fun provideLessonDao(
        fluentlyDatabase: FluentlyDatabase
    ): LessonDao = fluentlyDatabase.lessonDao()
}