package ru.fluentlyapp.fluently.database.entities

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "lessons")
data class LessonEntity(
    @PrimaryKey
    @ColumnInfo(name = "lesson_id")
    val lessonId: String,

    @ColumnInfo(name = "lesson_json")
    val lessonJson: String
)