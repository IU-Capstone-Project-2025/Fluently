package ru.fluentlyapp.fluently.feature.wordoftheday

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.datastore.WordOfTheDayDatastore
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepository
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgress
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepository
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import ru.fluentlyapp.fluently.network.model.WordOfTheDay
import timber.log.Timber
import java.time.Duration
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class WordOfTheDayRepository @Inject constructor(
    val fluentlyApiDataSource: FluentlyApiDataSource,
    val wordCacheRepository: WordCacheRepository,
    val wordProgressRepository: WordProgressRepository,
    val wordOfTheDayDatastore: WordOfTheDayDatastore
) {
    suspend fun updateWordOfTheDay() {
        val current = Instant.now()
        val lastSaved = wordOfTheDayDatastore.getLastSaved().first()
        if (lastSaved != null && Duration.between(lastSaved, current).toHours() < 24) {
            Timber.d("updateWordOfTheDay: the day hasn't passed since the last update")
            return
        }

        val result = fluentlyApiDataSource.getWordOfTheDay()
        wordOfTheDayDatastore.setWordOfTheDay(result)
        Timber.d("updateWordOfTheDay: $result")
    }

    fun getWordOfTheDay(): Flow<WordOfTheDay?> {
        return wordOfTheDayDatastore.getSavedWordOfTheDay()
    }

    suspend fun startLearningWordOfTheDay() {
        val wordOfTheDay = getWordOfTheDay().first()
        if (wordOfTheDay == null) {
            return
        }

        if (isWordOfTheDayLearning().first()) {
            Timber.d("startLearningWordOfTheDay: word is already learning")
            return
        }

        val wordCache = WordCache(
            wordId = wordOfTheDay.wordId,
            word = wordOfTheDay.word,
            translation = wordOfTheDay.translation,
            examples = wordOfTheDay.examples
        )
        val wordProgress = WordProgress(
            wordId = wordOfTheDay.wordId,
            isLearning = true,
            instant = Instant.now()
        )
        wordCacheRepository.updateWord(wordCache)
        wordProgressRepository.addProgress(wordProgress)
        Timber.d("Word $wordProgress is now learning")
    }

    fun isWordOfTheDayLearning(): Flow<Boolean> {
        val today = LocalDate.now()
        return wordProgressRepository.getProgresses(
            beginDate = today.atStartOfDay(ZoneId.systemDefault()).toInstant(),
            endDate = today.plusDays(1).atStartOfDay(ZoneId.systemDefault()).toInstant()
        ).map { list ->
            val wordOfTheDay = getWordOfTheDay().first()
            if (wordOfTheDay == null) {
                return@map false
            }
            list.find { it.wordId == wordOfTheDay.wordId } != null
        }
    }
}
