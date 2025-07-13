package ru.fluentlyapp.fluently.network

import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson
import ru.fluentlyapp.fluently.common.model.LessonComponent
import ru.fluentlyapp.fluently.network.model.internal.ExerciseApiModel.ExerciseType
import ru.fluentlyapp.fluently.network.model.internal.LessonResponseBody
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
                    }
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
                            selectedVariant = null
                        )
                    )
                }

                // Fill the gaps in the sentence exercise
//                ExerciseType.PICK_OPTION_SENTENCE -> {
//                    val options = exerciseData.pick_options!!.toMutableList()
//                    val correctVariant = options.insertRandomly(exerciseData.correct_answer!!)
//
//                    add(
//                        Exercise.FillTheGap(
//                            sentence = exerciseData.template!!.split("_".toRegex()),
//                            answerVariants = options,
//                            correctVariant = correctVariant,
//                            selectedVariant = null
//                        )
//                    )
//                }

                // Write the word from translation exercise
                ExerciseType.WRITE_WORD_FROM_TRANSLATION -> {
                    val translation = exerciseData.translation!!
                    val correctAnswer = exerciseData.correct_answer!!

                    add(
                        Exercise.InputWord(
                            translation = translation,
                            correctAnswer = correctAnswer,
                            inputtedWord = null
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