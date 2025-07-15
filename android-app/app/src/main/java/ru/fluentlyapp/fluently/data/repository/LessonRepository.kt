package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.distinctUntilChanged
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Exercise
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

class DefaultLessonRepository @Inject constructor(
    val fluentlyApiDataSource: FluentlyApiDataSource,
    val ongoingLessonDataStore: OngoingLessonDataStore,
    val wordCacheRepository: WordCacheRepository,
    val wordProgressRepository: WordProgressRepository
) : LessonRepository {
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
            add(Decoration.Finish())
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

    override suspend fun moveToNextComponent() {
        try {
            ongoingLessonDataStore.getOngoingLesson().first()?.let { lesson ->
                if (
                    lesson.currentLessonComponentIndex + 1 < lesson.components.size && // End not reached
                    (lesson.currentComponent as? Exercise)?.isAnswered != false // i.e. null (currentComponent is not exercise) or true
                ) {
                    ongoingLessonDataStore.setOngoingLesson(
                        lesson.copy(
                            currentLessonComponentIndex = lesson.currentLessonComponentIndex + 1
                        )
                    )
                }
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
        val progressMap = mutableMapOf<String, SentWordProgress>() // (word_id; progress)
        ongoingLessonDataStore.getOngoingLesson().first()?.let { lesson ->
            for (component in lesson.components) {
                if (component is Exercise.NewWord) {
                    var correctExercises = 0
                    var incorrectExercises = 0
                    for (possibleRelatedComponent in lesson.components) {
                        if (
                            possibleRelatedComponent is Exercise
                        ) {
                            val result = isExerciseCorrect(component, possibleRelatedComponent)
                            correctExercises += (result == true).compareTo(false)
                            incorrectExercises += (result == false).compareTo(false)
                        }
                    }
                    val overallExerciseCount = correctExercises + incorrectExercises
                    val correctnessRate: Float = correctExercises/ overallExerciseCount.toFloat()
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