package ru.fluentlyapp.fluently.database.app.wordprogress

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import java.time.Instant

@Entity(tableName = "word_progresses")
data class WordProgressEntity(
    @PrimaryKey val id: String,
    @ColumnInfo(name = "is_learning") val isLearning: Boolean,
    @ColumnInfo(name = "timestamp") val timestamp: Instant
)