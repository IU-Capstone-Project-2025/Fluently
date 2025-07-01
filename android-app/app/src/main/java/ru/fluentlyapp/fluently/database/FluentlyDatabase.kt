package ru.fluentlyapp.fluently.database

import androidx.room.Database
import androidx.room.RoomDatabase
import ru.fluentlyapp.fluently.database.dao.LessonDao
import ru.fluentlyapp.fluently.database.entities.LessonEntity

@Database(entities = [LessonEntity::class], version = 1)
abstract class FluentlyDatabase : RoomDatabase() {
    abstract fun lessonDao(): LessonDao
}