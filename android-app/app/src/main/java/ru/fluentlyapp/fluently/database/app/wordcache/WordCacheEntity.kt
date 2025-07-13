package ru.fluentlyapp.fluently.database.app.wordcache

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "word_caches")
data class WordCacheEntity(
    @PrimaryKey val id: String,
    @ColumnInfo(name = "word_json") val wordJson: String
)