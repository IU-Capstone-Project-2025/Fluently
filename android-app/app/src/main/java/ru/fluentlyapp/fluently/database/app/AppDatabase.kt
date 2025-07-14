package ru.fluentlyapp.fluently.database.app

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import ru.fluentlyapp.fluently.database.app.joinedwordprogress.JoinedWordProgressDao
import ru.fluentlyapp.fluently.database.app.wordcache.WordCacheEntity
import ru.fluentlyapp.fluently.database.app.util.InstantConverter
import ru.fluentlyapp.fluently.database.app.wordcache.WordCacheDao
import ru.fluentlyapp.fluently.database.app.wordprogress.WordProgressDao
import ru.fluentlyapp.fluently.database.app.wordprogress.WordProgressEntity

@Database(
    entities = [
        WordProgressEntity::class,
        WordCacheEntity::class
    ],
    version = 1,
)
@TypeConverters(InstantConverter::class)
abstract class AppDatabase : RoomDatabase() {
    abstract fun wordProgressDao(): WordProgressDao
    abstract fun wordCacheDao(): WordCacheDao
    abstract fun joinedWordProgressDao(): JoinedWordProgressDao
}