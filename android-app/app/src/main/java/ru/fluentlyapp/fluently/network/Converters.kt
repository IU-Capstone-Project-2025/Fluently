package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.network.model.Author
import ru.fluentlyapp.fluently.network.model.Chat
import ru.fluentlyapp.fluently.network.model.Message
import ru.fluentlyapp.fluently.network.model.Progress
import ru.fluentlyapp.fluently.network.model.UserPreferences
import ru.fluentlyapp.fluently.network.model.WordOfTheDay
import ru.fluentlyapp.fluently.network.model.internal.CardApiModel
import ru.fluentlyapp.fluently.network.model.internal.ChatRequestBody
import ru.fluentlyapp.fluently.network.model.internal.ChatResponseBody
import ru.fluentlyapp.fluently.network.model.internal.ExerciseApiModel.ExerciseType
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.internal.MessageApiModel
import ru.fluentlyapp.fluently.network.model.internal.UserPreferencesResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordOfTheDayResponseBody
import ru.fluentlyapp.fluently.network.model.internal.WordProgressApiModel
import kotlin.random.Random
import kotlin.random.nextInt


private fun <T> MutableList<T>.insertRandomly(item: T): Int {
    val insertPosition = Random.nextInt(0..size)
    add(insertPosition, item)
    return insertPosition
}

fun LessonResponseBody.convertToLesson(): Lesson {
    val lessonComponents = buildList<LessonComponent> {
        for (card in cards) {
            add(
                Exercise.NewWord(
                    word = card.word,
                    translation = card.translation,
                    phoneticTranscription = "",
                    doesUserKnow = null,
                    examples = card.sentences.map {
                        it.text to it.translation
                    },
                    wordId = card.word_id
                )
            )

            val exerciseData = card.exercise.data
            when (card.exercise.type) {
                // Pick the correct translation of the word exercise
                in listOf(
                    ExerciseType.TRANSLATE_EN_TO_RU,
                    ExerciseType.TRANSLATE_RU_TO_EN
                ) -> {
                    val options = exerciseData.pick_options!!.toMutableList()

                    add(
                        Exercise.ChooseTranslation(
                            word = exerciseData.text!!,
                            answerVariants = options,
                            correctVariant = options.indexOf(exerciseData.correct_answer),
                            selectedVariant = null,
                            wordId = card.word_id
                        )
                    )
                }

                // Fill the gaps in the sentence exercise
                ExerciseType.PICK_OPTION_SENTENCE -> {
                    val options = exerciseData.pick_options!!.toMutableList()
                    val expandedTemplate = " " + exerciseData.template + " "
                    add(
                        Exercise.FillTheGap(
                            sentence = expandedTemplate.split("_+".toRegex()),
                            answerVariants = options,
                            correctVariant = options.indexOf(exerciseData.correct_answer),
                            selectedVariant = null,
                            wordId = card.word_id
                        )
                    )
                }

                // Write the word from translation exercise
                ExerciseType.WRITE_WORD_FROM_TRANSLATION -> {
                    val translation = exerciseData.translation!!
                    val correctAnswer = exerciseData.correct_answer!!

                    add(
                        Exercise.InputWord(
                            translation = translation,
                            correctAnswer = correctAnswer,
                            inputtedWord = null,
                            wordId = card.word_id
                        )
                    )
                }
            }
        }
    }
    return Lesson(
        lessonId = "",
        components = lessonComponents,
        currentLessonComponentIndex = 0
    )
}

fun Progress.toProgressRequestBody() =
    progresses.map {
        WordProgressApiModel(
            cnt_reviewed = it.cntReviewed,
            confidence_score = it.confidenceScore,
            learned_at = it.learnedAt.toString(),
            word_id = it.wordId
        )
    }

fun WordOfTheDayResponseBody.toWordOfTheDay() = WordOfTheDay(
    wordId = word_id,
    word = word,
    translation = translation,
    examples = sentences.map {
        it.text to it.translation
    }
)

fun Message.toMessageApiModel() = MessageApiModel(
    author = author.key,
    message = message
)

fun MessageApiModel.toMessage() = Message(
    author = Author.entries.first { it.key == author },
    message = message
)

fun ChatResponseBody.toChat() = Chat(
    chat = chat.map { it.toMessage() }
)

fun Chat.toChatRequestBody() = ChatRequestBody(
    chat = chat.map { it.toMessageApiModel() }
)


fun UserPreferencesResponseBody.toUserPreferences(): UserPreferences {
    return UserPreferences(
        avatarImageUrl = avatar_image_url,
        cefrLevel = cefr_level,
        factEveryday = fact_everyday,
        goal = goal,
        id = id,
        notificationAt = notification_at,
        notifications = notifications,
        subscribed = subscribed,
        userId = user_id,
        wordsPerDay = words_per_day
    )
}
