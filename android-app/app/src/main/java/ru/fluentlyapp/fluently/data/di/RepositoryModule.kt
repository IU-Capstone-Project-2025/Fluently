package ru.fluentlyapp.fluently.data.di

import dagger.Binds
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import ru.fluentlyapp.fluently.data.repository.AuthRepository
import ru.fluentlyapp.fluently.data.repository.GoogleBasedAuthRepository
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.data.repository.StubLessonRepository
import ru.fluentlyapp.fluently.oauth.GoogleOAuthService
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class RepositoryModule {
    @Binds
    @Singleton
    abstract fun bindLessonRepository(
        stubLessonRepository: StubLessonRepository
    ): LessonRepository

    @Binds
    @Singleton
    abstract fun bindAuthRepository(
        googleBasedAuthRepository: GoogleBasedAuthRepository
    ): AuthRepository
}