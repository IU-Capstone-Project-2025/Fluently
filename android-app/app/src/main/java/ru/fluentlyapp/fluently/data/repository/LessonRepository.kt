package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.update
import ru.fluentlyapp.fluently.datastore.LessonPreferencesDataStore
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.Lesson
import ru.fluentlyapp.fluently.network.FluentlyDataSource
import ru.fluentlyapp.fluently.network.di.mockLessonResponse
import javax.inject.Inject

interface LessonRepository {
    /**
     * Return the id of locally set ongoing lesson.
     *
     * Returns null if none of the lessons are ongoing.
     */
    fun getSavedOngoingLessonIdAsFlow(): Flow<String?>

    /**
     * Locally, set the `lessonId` as the ongoing lesson.
     */
    suspend fun setSavedOngoingLessonId(lessonId: String)

    /**
     * Locally, drop the ongoing lesson id.
     */
    suspend fun dropSavedOngoingLesson()

    /**
     * Get the saved lesson as `Flow` by the `lessonId`. The flow may emit null if
     * none of the lessons are stored under `lessonId`.
     */
    fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?>

    /**
     * At any moment, any user has the current ongoing lesson. This method fetches
     * the currently assigned lesson for this user.
     *
     * May throw exception.
     */
    suspend fun fetchCurrentLesson(): Lesson

    /**
     * Fetch the lesson by the `lessonId` from the server
     *
     * May throw exception.
     */
    suspend fun fetchLesson(lessonId: String): Lesson

    /**
     * Update the `lesson` locally.
     */
    suspend fun saveLesson(lesson: Lesson)

    /**
     * Get the saved lesson by the `lessonId`.
     *
     * Returns null if no lessons are saved under `lessonId`.
     */
    suspend fun getSavedLesson(lessonId: String): Lesson?

    /**
     * Send the lesson to the server so that it stores it.
     *
     * May throw exception.
     */
    suspend fun sendLesson(lesson: Lesson)
}

var testLesson = Lesson(
    lessonId = "lesson_test_001",
    components = listOf(
        // Word 1: Consciousness
        Exercise.NewWord(
            word = "consciousness",
            translation = "сознание",
            phoneticTranscription = "/ˈkɑːn.ʃəs.nəs/",
            doesUserKnow = null,
            examples = listOf("She lost consciousness" to "Она потеряла сознание")
        ),
        Exercise.ChooseTranslation(
            word = "consciousness",
            answerVariants = listOf("сознание", "память", "внимание", "мысль"),
            correctVariant = 0,
            selectedVariant = null
        ),

        // Word 2: Awareness
        Exercise.NewWord(
            word = "awareness",
            translation = "осознание",
            phoneticTranscription = "/əˈweə.nəs/",
            doesUserKnow = null,
            examples = listOf("Environmental awareness is rising" to "Экологическое осознание растет")
        ),
        Exercise.ChooseTranslation(
            word = "осознание",
            answerVariants = listOf("consciousness", "awareness", "focus", "clarity"),
            correctVariant = 1,
            selectedVariant = null
        ),

        // Word 3: Resilience
        Exercise.NewWord(
            word = "resilience",
            translation = "устойчивость",
            phoneticTranscription = "/rɪˈzɪl.jəns/",
            doesUserKnow = null,
            examples = listOf("Resilience is key to recovery" to "Устойчивость — ключ к восстановлению")
        ),
        Exercise.ChooseTranslation(
            word = "resilience",
            answerVariants = listOf("гибкость", "восстановление", "устойчивость", "усилие"),
            correctVariant = 2,
            selectedVariant = null
        ),

        // Word 4: Determination
        Exercise.NewWord(
            word = "determination",
            translation = "решимость",
            phoneticTranscription = "/dɪˌtɜː.mɪˈneɪ.ʃən/",
            doesUserKnow = null,
            examples = listOf("Her determination inspired others" to "Её решимость вдохновляла других")
        ),
        Exercise.ChooseTranslation(
            word = "решимость",
            answerVariants = listOf("motivation", "decision", "goal", "determination"),
            correctVariant = 3,
            selectedVariant = null
        ),

        // Word 5: Empathy
        Exercise.NewWord(
            word = "empathy",
            translation = "сочувствие",
            phoneticTranscription = "/ˈem.pə.θi/",
            doesUserKnow = null,
            examples = listOf("Empathy helps build trust" to "Сочувствие помогает строить доверие")
        ),
        Exercise.ChooseTranslation(
            word = "empathy",
            answerVariants = listOf("sympathy", "compassion", "empathy", "pity"),
            correctVariant = 2,
            selectedVariant = null
        )
    )
)


class StubLessonRepository @Inject constructor(
    val lessonPreferencesDataStore: LessonPreferencesDataStore,
    val fluentlyDataSource: FluentlyDataSource
) : LessonRepository {
    val lessons = mutableMapOf<String, MutableStateFlow<Lesson?>>()
    override fun getSavedOngoingLessonIdAsFlow(): Flow<String?> {
        return lessonPreferencesDataStore.getOngoingLessonIdAsFlow()
    }

    override suspend fun setSavedOngoingLessonId(lessonId: String) {
        lessonPreferencesDataStore.setOngoingLessonId(lessonId)
    }

    override suspend fun dropSavedOngoingLesson() {
        lessonPreferencesDataStore.dropOngoingLessonId()
    }

    override suspend fun fetchCurrentLesson(): Lesson {
        return fluentlyDataSource.getCurrentLesson()
    }

    override suspend fun fetchLesson(lessonId: String): Lesson {
        if (lessonId == mockLessonResponse.lesson.lesson_id) {
            return fluentlyDataSource.getCurrentLesson()
        } else {
            TODO("Not yet implemented")
        }
    }

    override suspend fun sendLesson(lesson: Lesson) {
        TODO("Not yet implemented")
    }

    override suspend fun saveLesson(lesson: Lesson) {
        val stateFlow = lessons.getOrPut(lesson.lessonId) {
            MutableStateFlow(null)
        }
        stateFlow.update { lesson }
    }

    override suspend fun getSavedLesson(lessonId: String): Lesson? {
        return lessons[lessonId]?.value
    }

    override fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?> {
        val stateFlow = lessons.getOrPut(lessonId) {
            MutableStateFlow(null)
        }

        return stateFlow
    }
}