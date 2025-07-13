package ru.fluentlyapp.fluently.feature.exercisecache.database

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

enum class ExerciseType {
    NEW_WORD, CHOOSE_TRANSLATION, WRITE_WORD, FILL_THE_GAP
}

@Entity(tableName = "exercise_cache")
data class ExerciseCacheEntity(
    @PrimaryKey val id: String,
    @ColumnInfo(name = "exercise_type") val exerciseType: ExerciseType,
    @ColumnInfo(name = "exercise_json") val exerciseJson: String
)