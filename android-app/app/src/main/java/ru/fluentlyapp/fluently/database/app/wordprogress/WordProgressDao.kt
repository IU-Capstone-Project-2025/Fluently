package ru.fluentlyapp.fluently.database.app.wordprogress

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import java.time.Instant

@Dao
interface WordProgressDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(wordProgress: WordProgressEntity)

    @Query(
        "SELECT * FROM word_progresses WHERE timestamp BETWEEN :begin AND :end"
    )
    fun getProgressesBetweenDates(
        begin: Instant,
        end: Instant
    ): Flow<List<WordProgressEntity>>
}