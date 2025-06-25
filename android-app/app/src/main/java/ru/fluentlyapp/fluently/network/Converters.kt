package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.data.model.ServerToken
import ru.fluentlyapp.fluently.model.Exercise
import ru.fluentlyapp.fluently.model.Lesson
import ru.fluentlyapp.fluently.model.LessonComponent
import ru.fluentlyapp.fluently.network.model.ExerciseApiModel.ExerciseType
import ru.fluentlyapp.fluently.network.model.LessonResponseBody
import ru.fluentlyapp.fluently.network.model.ServerTokenResponseBody
import kotlin.random.Random
import kotlin.random.nextInt

fun ServerTokenResponseBody.toServerToken() = ServerToken(
    accessToken = accessToken,
    refreshToken = refreshToken,
    tokenType = tokenType,
    expiresInSeconds = expiresInSeconds
)

private fun <T> MutableList<T>.insertRandomly(item: T): Int {
    val insertPosition = Random.nextInt(0..size)
    add(insertPosition, item)
    return insertPosition
}

fun LessonResponseBody.convertToLesson(): Lesson {
    val lessonComponents = buildList<LessonComponent> {
        for (card in cards) {
            // First, add the word
            if (card.is_new) {
                add(
                    Exercise.NewWord(
                        word = card.word,
                        translation = card.translation,
                        phoneticTranscription = card.translation,
                        doesUserKnow = null,
                        examples = card.sentences.map {
                            it.text to it.translation
                        }
                    )
                )
            }

            // Second, add the related exercise
            val exerciseData = card.exercise.data
            when {
                // Translation and pick option exercise
                card.exercise.type in listOf(
                    ExerciseType.TRANSLATE_EN_TO_RU,
                    ExerciseType.TRANSLATE_RU_TO_EN
                ) -> {
                    val options = exerciseData.pick_options!!.toMutableList()
                    val correctVariant = options.insertRandomly(exerciseData.correct_answer!!)

                    add(
                        Exercise.ChooseTranslation(
                            word = card.word,
                            answerVariants = options,
                            correctVariant = correctVariant,
                            selectedVariant = null
                        )
                    )
                }

                // Fill the gaps in sentence exercise
                card.exercise.type == ExerciseType.PICK_OPTIONS_SENTENCE -> {
                    val options = exerciseData.pick_options!!.toMutableList()
                    val correctVariant = options.insertRandomly(exerciseData.correct_answer!!)

                    add(
                        Exercise.FillTheGap(
                            sentence = exerciseData.template!!.split("\\s+".toRegex()),
                            answerVariants = options,
                            correctVariant = correctVariant,
                            selectedVariant = null
                        )
                    )
                }
            }
        }
    }
    return Lesson(
        lessonId = lesson.lesson_id,
        components = lessonComponents,
        currentLessonComponentIndex = 0
    )
}