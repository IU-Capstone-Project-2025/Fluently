package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.distinctUntilChanged
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.datastore.OngoingLessonDataStore
import ru.fluentlyapp.fluently.feature.wordcache.WordCache
import ru.fluentlyapp.fluently.feature.wordcache.WordCacheRepository
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgress
import ru.fluentlyapp.fluently.feature.wordprogress.WordProgressRepository
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import ru.fluentlyapp.fluently.network.model.Progress
import ru.fluentlyapp.fluently.network.model.SentWordProgress
import timber.log.Timber
import java.time.Instant
import javax.inject.Inject

data class LessonStatistic(
    val learnedWordsCount: Int = 0,
    val wordsInProgress: Int = 0
)

interface LessonRepository {
    /**
     * Helper method to quickly determine if there is a saved ongoing lesson.
     */
    suspend fun hasSavedLesson(): Boolean

    /**
     * Ask the server to generate new lesson for the user, then fetch it and store it locally.
     *
     * May throw exception.
     */
    suspend fun fetchAndSaveOngoingLesson()


    /**
     * Try to update the current lesson component.
     *
     * If the update cannot happen (the type of `updatedComponent` is not the same as
     * the type the ongoing lesson is on, or the value of updatedComponent cannot be applied for
     * some other reason), the method will have no effect.
     *
     * The method may throw exception.
     */
    suspend fun updateCurrentComponent(updatedComponent: LessonComponent)

    /**
     * Get the current `LessonComponent` as the `Flow`. If there is no saved `ongoingLesson`, the
     * flow may emit null.
     *
     * During the collection, the returned `Flow` may throw exceptions.
     */
    fun currentComponent(): Flow<LessonComponent?>

    /**
     * If it is possible (e.g. the exercise is completed), moves to the next lesson component.
     * If the end is reached, the method have no effect.
     */
    suspend fun moveToNextComponent()

    /**
     * Once the user finishes the lesson, send the progress to the api.
     */
    suspend fun finishLesson()
}

const val PREFERRED_NUMBER_OF_WORDS = 10

class DefaultLessonRepository @Inject constructor(
    val fluentlyApiDataSource: FluentlyApiDataSource,
    val ongoingLessonDataStore: OngoingLessonDataStore,
    val wordCacheRepository: WordCacheRepository,
    val wordProgressRepository: WordProgressRepository
) : LessonRepository {
    private enum class NewWordExerciseStatus {
        IGNORED,
        USER_IS_LEARNING,
        USER_KNOWS,
        NO_OCCURRENCE
    }

    override suspend fun hasSavedLesson(): Boolean {
        return try {
            ongoingLessonDataStore.getOngoingLesson().first() != null
        } catch (ex: Exception) {
            Timber.e(ex)
            false
        }
    }

    private fun generateOnboardingComponent(components: List<LessonComponent>): Decoration.Onboarding {
        var wordsCount = 0
        var exercisesCount = 0
        for (component in components) {
            if (component is Exercise.NewWord) {
                wordsCount++
            } else if (component is Exercise) {
                exercisesCount++
            }
        }
        return Decoration.Onboarding(wordsCount, exercisesCount)
    }

    private fun List<LessonComponent>.withIdSetToIndex(): List<LessonComponent> {
        for (index in indices) {
            this[index].id = index
        }
        return this
    }

    override suspend fun fetchAndSaveOngoingLesson() {
        val lesson = fluentlyApiDataSource.getLesson()

        for (component in lesson.components) {
            if (component is Exercise.NewWord) {
                wordCacheRepository.updateWord(
                    WordCache(
                        wordId = component.wordId,
                        word = component.word,
                        translation = component.translation,
                        examples = component.examples
                    )
                )
            }
        }

        val updatedLessonComponents: List<LessonComponent> = buildList {
            add(generateOnboardingComponent(lesson.components))
            addAll(lesson.components)
        }.withIdSetToIndex()

        ongoingLessonDataStore.setOngoingLesson(lesson.copy(components = updatedLessonComponents))
        Timber.d("Store the received lesson")
    }

    override suspend fun updateCurrentComponent(updatedComponent: LessonComponent) {
        ongoingLessonDataStore.getOngoingLesson().first()?.let { lesson ->
            if (
                updatedComponent::class == lesson.currentComponent::class &&
                updatedComponent is Exercise &&
                (lesson.currentComponent as? Exercise)?.isAnswered == false &&
                updatedComponent.isAnswered
            ) {
                Timber.d("The `updatedComponent`=$updatedComponent is valid for update")
                val updatedComponents = lesson.components.toMutableList().apply {
                    this[lesson.currentLessonComponentIndex] = updatedComponent
                }
                val updatedLesson = lesson.copy(
                    components = updatedComponents
                )
                ongoingLessonDataStore.setOngoingLesson(updatedLesson)
                Timber.d("Save the updated lesson: $lesson")
            }
            Timber.d("The `updatedComponent`=$updatedComponent is NOT valid for update - ignore")
        }
    }

    private fun validateCandidateLessonComponent(
        lesson: Lesson,
        candidateComponentIndex: Int
    ): Boolean {
        val candidateComponent = lesson.components[candidateComponentIndex]
        val currentWordId = when (candidateComponent) {
            is Exercise.FillTheGap -> { candidateComponent.wordId }
            is Exercise.InputWord -> {  candidateComponent.wordId }
            is Exercise.ChooseTranslation -> { candidateComponent.wordId }
            else -> null
        }

        /**
         * The rules are the following:
         * 1) If component candidate is new word, then show it only if the user hasn't already
         * started to learn a certain threshold of words
         * 2) If component candidate is some exercise related to some word, then show it only if
         * the related word is either hasn't been previously met in the lesson or the set that they
         * do not know it
         */
        var learningWordsCount = 0
        var originalWordStatus: NewWordExerciseStatus = NewWordExerciseStatus.NO_OCCURRENCE
        for (i in 0..<candidateComponentIndex) {
            val component = lesson.components[i]
            if (component is Exercise.NewWord) {
                if (component.doesUserKnow == false) {
                    learningWordsCount++
                }
                if (currentWordId != null && component.wordId == currentWordId) {
                    originalWordStatus = when (component.doesUserKnow) {
                        true -> NewWordExerciseStatus.USER_KNOWS
                        false -> NewWordExerciseStatus.USER_IS_LEARNING
                        null -> NewWordExerciseStatus.IGNORED
                    }
                }
            }
        }
        Timber.d("learningWordsCount=$learningWordsCount; originalWordStatus=$originalWordStatus")
        return if (candidateComponent is Exercise.NewWord) {
            learningWordsCount < PREFERRED_NUMBER_OF_WORDS
        } else if (currentWordId != null) {
            originalWordStatus in setOf(NewWordExerciseStatus.NO_OCCURRENCE, NewWordExerciseStatus.USER_IS_LEARNING)
        } else {
            true
        }
    }

    override suspend fun moveToNextComponent() {
        try {
            val lesson = ongoingLessonDataStore.getOngoingLesson().first()
            if (lesson == null) {
                return
            }

            if (
                (lesson.currentComponent as? Exercise)?.isAnswered == false // current exercise is not answered
            ) {
                return
            }

            if (lesson.currentLessonComponentIndex + 1 == lesson.components.size) {
                // Add the finish screen
                val updatedLessonComponents: List<LessonComponent> = buildList {
                    addAll(lesson.components)
                    add(Decoration.Finish())
                }.withIdSetToIndex()
                ongoingLessonDataStore.setOngoingLesson(
                    lesson.copy(
                        currentLessonComponentIndex = lesson.currentLessonComponentIndex + 1,
                        components = updatedLessonComponents
                    )
                )
            }

            var candidateComponentIndex = lesson.currentLessonComponentIndex + 1
            while (
                candidateComponentIndex < lesson.components.size &&
                !validateCandidateLessonComponent(lesson, candidateComponentIndex)
            ) {
                candidateComponentIndex++
            }

            if (candidateComponentIndex < lesson.components.size) {
                ongoingLessonDataStore.setOngoingLesson(
                    lesson.copy(
                        currentLessonComponentIndex = candidateComponentIndex
                    )
                )
            } else {
                val updatedLessonComponents: List<LessonComponent> = buildList {
                    addAll(lesson.components)
                    add(Decoration.Finish())
                }.withIdSetToIndex()
                ongoingLessonDataStore.setOngoingLesson(
                    lesson.copy(
                        currentLessonComponentIndex = lesson.components.size,
                        components = updatedLessonComponents
                    )
                )
            }
        } catch (ex: Exception) {
            Timber.e(ex)
        }
    }

    override fun currentComponent(): Flow<LessonComponent?> {
        return ongoingLessonDataStore
            .getOngoingLesson()
            .distinctUntilChanged()
            .map { lesson ->
                val component = lesson?.currentComponent
                Timber.v("Emit lessonComponent: $component")
                component
            }
    }

    private fun isExerciseCorrect(word: Exercise.NewWord, exercise: Exercise): Boolean? {
        return if (
            (exercise is Exercise.FillTheGap && exercise.wordId == word.wordId) ||
            (exercise is Exercise.NewWord && exercise.wordId == word.wordId) ||
            (exercise is Exercise.InputWord && exercise.wordId == word.wordId) ||
            (exercise is Exercise.ChooseTranslation && exercise.wordId == word.wordId)
        )
            exercise.isCorrect
        else
            null

    }

    override suspend fun finishLesson() {
        // Consider only words that HAS been answered
        val progressMap = mutableMapOf<String, SentWordProgress>() // (word_id; progress)
        ongoingLessonDataStore.getOngoingLesson().first()?.let { lesson ->
            for (component in lesson.components) {
                if (component is Exercise.NewWord && component.isAnswered) {
                    var correctExercises = 0
                    var incorrectExercises = 0
                    for (possiblyRelatedComponent in lesson.components) {
                        if (
                            possiblyRelatedComponent is Exercise
                        ) {
                            val result = isExerciseCorrect(component, possiblyRelatedComponent)
                            correctExercises += (result == true).compareTo(false)
                            incorrectExercises += (result == false).compareTo(false)
                        }
                    }
                    val overallExerciseCount = correctExercises + incorrectExercises
                    val correctnessRate: Float = correctExercises / overallExerciseCount.toFloat()
                    val progress = SentWordProgress(
                        wordId = component.wordId,
                        cntReviewed = overallExerciseCount,
                        confidenceScore = (correctnessRate * 100).toInt().coerceIn(0, 100),
                        learnedAt = Instant.now()
                    )
                    progressMap[component.wordId] = progress
                }
            }
        }
        Timber.v("Compose the progress map: ${progressMap.values.joinToString(", ")}")
        for (value in progressMap.values) {
            wordProgressRepository.addProgress(
                WordProgress(
                    wordId = value.wordId,
                    isLearning = value.confidenceScore < 90,
                    instant = value.learnedAt
                )
            )
        }
        Timber.v("Save to the progress map")

        fluentlyApiDataSource.sendProgress(
            Progress(
                progresses = progressMap.values.toList()
            )
        )
        Timber.v("Send to the fluently api data source")

        ongoingLessonDataStore.dropOngoingLesson()
        Timber.d("Drop the ongoing lesson")
    }
}