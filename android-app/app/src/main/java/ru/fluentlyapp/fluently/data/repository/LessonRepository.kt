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
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import timber.log.Timber
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
     * Clear the ongoing lesson content.
     *
     * May throw exception.
     */
    suspend fun dropOngoingLesson()

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
    suspend fun sendLesson()
}

class DefaultLessonRepository @Inject constructor(
    val fluentlyApiDataSource: FluentlyApiDataSource,
    val ongoingLessonDataStore: OngoingLessonDataStore,
    val wordCacheRepository: WordCacheRepository
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

    override suspend fun dropOngoingLesson() {
        ongoingLessonDataStore.dropOngoingLesson()
        Timber.d("Drop the ongoing lesson")
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

    override suspend fun sendLesson() {
        ongoingLessonDataStore.getOngoingLesson().first()?.let { lesson ->

        }
    }
}