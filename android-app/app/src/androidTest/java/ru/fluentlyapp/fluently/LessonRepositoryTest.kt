package ru.fluentlyapp.fluently

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.preferencesDataStore
import androidx.test.platform.app.InstrumentationRegistry
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import org.junit.Before
import org.junit.Test
import ru.fluentlyapp.fluently.common.model.Decoration
import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.data.repository.DefaultLessonRepository
import ru.fluentlyapp.fluently.data.repository.LessonRepository
import ru.fluentlyapp.fluently.datastore.OngoingLessonDataStore
import ru.fluentlyapp.fluently.network.FluentlyApiDataSource
import org.junit.Assert.*
import ru.fluentlyapp.fluently.common.model.LessonComponent
import timber.log.Timber

private val testLesson = Lesson(
    lessonId = "test_lesson_001",
    components = listOf(
        Exercise.NewWord(
            word = "Consciousness",
            translation = "Сознание",
            phoneticTranscription = "/ˈkɑːn.ʃəs.nəs/",
            doesUserKnow = null,
            examples = listOf(
                "Human consciousness is complex." to "Человеческое сознание сложно.",
                "She lost consciousness after the accident." to "Она потеряла сознание после аварии."
            )
        ),
        Exercise.ChooseTranslation(
            word = "Apple",
            answerVariants = listOf("Яблоко", "Апельсин", "Банан"),
            correctVariant = 0,
            selectedVariant = null
        ),
        Decoration.Loading
    ),
    currentLessonComponentIndex = 0
)

class DefaultLessonRepositoryTest {
    private val Context.testDataStore by preferencesDataStore(
        "test_datastore_${System.currentTimeMillis()}"
    )
    lateinit var context: Context
    lateinit var lessonRepository: LessonRepository

    suspend fun clearTestDataStore() {
        context.testDataStore.edit { it.clear() }
    }

    @Before
    fun init() {
        Timber.plant(Timber.DebugTree())

        context = InstrumentationRegistry.getInstrumentation().targetContext
        val ongoingLessonDataStore = OngoingLessonDataStore(context.testDataStore)
        val fluentlyApiDataSource = object : FluentlyApiDataSource {
            override suspend fun getLesson(): Lesson = testLesson
        }
        lessonRepository = DefaultLessonRepository(
            fluentlyApiDataSource,
            ongoingLessonDataStore
        )
    }

    @Test
    fun currentComponentAfterFetchAndSave_producesCorrectLessonComponent() {
        runBlocking {
            clearTestDataStore()

            // Fetch the ongoing lesson
            lessonRepository.fetchAndSaveOngoingLesson()
            val currentComponent: LessonComponent? = lessonRepository.currentComponent().first()
            assertEquals(testLesson.currentComponent, currentComponent)
        }
    }

    @Test
    fun dropOngoingLesson_producesNullInSubsequentFlows() {
        runBlocking {
            clearTestDataStore()

            lessonRepository.fetchAndSaveOngoingLesson()
            lessonRepository.dropOngoingLesson()
            val currentComponent = lessonRepository.currentComponent().first()
            assertEquals(null, currentComponent)
        }
    }

    @Test
    fun updateCurrentComponent_correctlyUpdatesAndHandlesTheComponent() {
        runBlocking {
            clearTestDataStore()

            lessonRepository.fetchAndSaveOngoingLesson()
            // Finally, answered exercise of the same type should update the current exercise
            val updatedComponent = (testLesson.components[0] as Exercise.NewWord).copy(
                doesUserKnow = true
            )
            lessonRepository.updateCurrentComponent(updatedComponent)
            val currentComponent = lessonRepository.currentComponent().first()
            assertEquals(updatedComponent, currentComponent)
        }
    }
}