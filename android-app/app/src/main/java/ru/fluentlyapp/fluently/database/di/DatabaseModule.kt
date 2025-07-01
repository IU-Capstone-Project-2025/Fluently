package ru.fluentlyapp.fluently.database.di

import android.content.Context
import androidx.room.Room
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.database.FluentlyDatabase
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
class DatabaseModule {
    @Provides
    @Singleton
    fun provideFluentlyDatabase(
        @ApplicationContext appContext: Context
    ): FluentlyDatabase {
        return Room.databaseBuilder(
            appContext,
            FluentlyDatabase::class.java,
            "fluently-database"
        ).build()
    }
}