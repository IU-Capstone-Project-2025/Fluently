package ru.fluentlyapp.fluently.data.repository

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.Lesson
import javax.inject.Inject

interface LessonRepository {
    /**
     * Return the id of locally set ongoing lesson.
     *
     * Returns null if none of the lessons are ongoing.
     */
    suspend fun getLocalOngoingLessonId(): String?

    /**
     * Locally, set the `lessonId` as the ongoing lesson.
     */
    suspend fun setLocalOngoingLessonId(lessonId: String)

    /**
     * Locally, mark the none of the lessons are ongoing.
     */
    suspend fun dropOngoingLesson()

    suspend fun fetchCurrentLesson(): Lesson
    suspend fun fetchLesson(lessonId: String): Lesson
    suspend fun saveLesson(lesson: Lesson)
    suspend fun getSavedLesson(lessonId: String): Lesson?
    suspend fun moveToNextComponent(lessonId: String)
    fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?>
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


class StubLessonRepository @Inject constructor() : LessonRepository {
    private val lesson = MutableStateFlow(testLesson)
    override suspend fun getLocalOngoingLessonId(): String? {
        TODO("Not yet implemented")
    }

    override suspend fun setLocalOngoingLessonId(lessonId: String) {
        TODO("Not yet implemented")
    }


    override suspend fun dropOngoingLesson() {
        TODO("Not yet implemented")
    }

    override suspend fun fetchCurrentLesson(): Lesson {
        TODO("Not yet implemented")
    }

    override suspend fun fetchLesson(lessonId: String): Lesson {
        TODO("Not yet implemented")
    }

    override suspend fun saveLesson(lesson: Lesson) {
        this.lesson.update { lesson }
    }

    override suspend fun getSavedLesson(lessonId: String): Lesson? {
        return lesson.value
    }

    override suspend fun moveToNextComponent(lessonId: String) {
        val newComponentIndex = with(lesson.value) {
            (currentLessonComponentIndex + 1).coerceIn(0, components.size - 1)
        }
        lesson.update { it.copy(currentLessonComponentIndex = newComponentIndex) }
    }

    override fun getSavedLessonAsFlow(lessonId: String): Flow<Lesson?> = lesson.asStateFlow()
}