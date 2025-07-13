package ru.fluentlyapp.fluently.feature.wordprogress.database

import androidx.room.TypeConverter
import java.time.Instant

class TimestampConverter {
    @TypeConverter
    fun instantToLong(instant: Instant): Long {
        return instant.toEpochMilli()
    }

    @TypeConverter
    fun epochMilliToInstant(epochMilli: Long): Instant {
        return Instant.ofEpochMilli(epochMilli)
    }
}