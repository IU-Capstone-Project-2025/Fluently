package ru.fluentlyapp.fluently.database.app.joinedwordprogress

import androidx.room.Dao
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import java.time.Instant

@Dao
interface JoinedWordProgressDao {
    @Query("""
        SELECT 
        wc.id AS id,
        wp.is_learning AS is_learning,
        wp.timestamp AS timestamp,
        wc.word_json AS word_json
        FROM
        word_caches AS wc
        INNER JOIN word_progresses AS wp 
        ON wc.id = wp.id
        WHERE wp.timestamp BETWEEN :begin AND :end 
    """)
    fun getJoinedWordProgress(
        begin: Instant,
        end: Instant
    ): Flow<List<JoinedWordProgressData>>
}