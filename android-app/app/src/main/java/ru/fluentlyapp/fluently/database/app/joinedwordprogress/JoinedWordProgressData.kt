package ru.fluentlyapp.fluently.database.app.joinedwordprogress

import androidx.room.ColumnInfo
import java.time.Instant

data class JoinedWordProgressData(
    val id: String,
    @ColumnInfo(name = "is_learning") val isLearning: Boolean,
    val timestamp: Instant,
    @ColumnInfo("word_json") val wordJson: String
)