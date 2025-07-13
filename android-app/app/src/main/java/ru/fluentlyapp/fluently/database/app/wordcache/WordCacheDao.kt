package ru.fluentlyapp.fluently.database.app.wordcache

import androidx.room.Dao
import androidx.room.Delete
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow

@Dao
interface WordCacheDao {
    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(wordCacheEntity: WordCacheEntity)

    @Delete
    suspend fun delete(wordCacheEntity: WordCacheEntity)

    @Query("SELECT * FROM word_caches")
    fun getAll(): Flow<List<WordCacheEntity>>

    @Query("SELECT * FROM word_caches WHERE id = :id")
    fun getById(id: String): Flow<WordCacheEntity>
}