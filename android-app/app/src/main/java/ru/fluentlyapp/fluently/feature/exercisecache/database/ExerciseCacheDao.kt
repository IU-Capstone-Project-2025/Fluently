package ru.fluentlyapp.fluently.feature.exercisecache.database

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query

@Dao
interface ExerciseCacheDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(exerciseCacheEntity: ExerciseCacheEntity)

    @Delete
    suspend fun delete(exerciseCacheEntity: ExerciseCacheEntity)

    @Query("SELECT * FROM exercise_cache WHERE id = :id")
    suspend fun getById(id: String): ExerciseCacheEntity
}
