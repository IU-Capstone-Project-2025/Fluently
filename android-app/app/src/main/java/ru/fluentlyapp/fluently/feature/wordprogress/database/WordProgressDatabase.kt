package ru.fluentlyapp.fluently.feature.wordprogress.database

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import java.sql.Timestamp

@Database(entities = [WordProgressEntity::class], version = 1)
@TypeConverters(TimestampConverter::class)
abstract class WordProgressDatabase : RoomDatabase() {
    abstract fun wordProgressDao(): WordProgressDao
}