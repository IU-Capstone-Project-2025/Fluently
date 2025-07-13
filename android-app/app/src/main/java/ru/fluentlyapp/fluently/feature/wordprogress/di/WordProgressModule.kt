package ru.fluentlyapp.fluently.feature.wordprogress.di

import android.content.Context
import androidx.room.Room
import dagger.Binds
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepository
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepositoryImpl
import ru.fluentlyapp.fluently.feature.wordprogress.database.WordProgressDatabase
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)

abstract class WordProgressModule {
    @Binds
    @Singleton
    abstract fun bindWordProgressRepository(
        wordProgressRepositoryImpl: WordProgressRepositoryImpl
    ): WordProgressRepository

    companion object {
        @Provides
        @Singleton
        fun provideWordProgressDatabase(
            @ApplicationContext context: Context
        ): WordProgressDatabase {
            return Room.databaseBuilder(
                context,
                WordProgressDatabase::class.java,
                "word-progress-database"
            ).build()
        }
    }
}