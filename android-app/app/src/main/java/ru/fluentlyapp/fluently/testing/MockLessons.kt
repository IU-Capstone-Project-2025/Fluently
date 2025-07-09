package ru.fluentlyapp.fluently.testing

import ru.fluentlyapp.fluently.common.model.Exercise
import ru.fluentlyapp.fluently.common.model.Lesson

object MockLessons {
    var lessonIndex = 0
    val lessons: List<Lesson> = listOf(
        Lesson(
            lessonId = "lesson_real_01",
            components = buildList {
                addAll(
                    listOf(
                        Exercise.NewWord(
                            word = "Confidence",
                            translation = "Уверенность",
                            phoneticTranscription = "/ˈkɒn.fɪ.dəns/",
                            doesUserKnow = null,
                            examples = listOf(
                                "She spoke with confidence." to "Она говорила с уверенностью.",
                                "Confidence is important in public speaking." to "Уверенность важна при публичных выступлениях."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Journey",
                            translation = "Путешествие",
                            phoneticTranscription = "/ˈdʒɜː.ni/",
                            doesUserKnow = null,
                            examples = listOf(
                                "The journey took five hours." to "Путешествие заняло пять часов.",
                                "It's about the journey, not the destination." to "Важно путешествие, а не пункт назначения."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Effort",
                            translation = "Усилие",
                            phoneticTranscription = "/ˈef.ət/",
                            doesUserKnow = null,
                            examples = listOf(
                                "It requires a lot of effort to succeed." to "Для успеха требуется много усилий.",
                                "He made an effort to improve his English." to "Он приложил усилия, чтобы улучшить свой английский."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Improve",
                            translation = "Улучшать",
                            phoneticTranscription = "/ɪmˈpruːv/",
                            doesUserKnow = null,
                            examples = listOf(
                                "Practice will improve your skills." to "Практика улучшит твои навыки.",
                                "I want to improve my pronunciation." to "Я хочу улучшить своё произношение."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Mistake",
                            translation = "Ошибка",
                            phoneticTranscription = "/mɪˈsteɪk/",
                            doesUserKnow = null,
                            examples = listOf(
                                "It's okay to make mistakes." to "Ошибаться — это нормально.",
                                "We learn from our mistakes." to "Мы учимся на своих ошибках."
                            )
                        )
                    )
                )

                addAll(
                    listOf(
                        Exercise.ChooseTranslation(
                            word = "Confidence",
                            answerVariants = listOf("Уверенность", "Сомнение", "Страх"),
                            correctVariant = 0,
                            selectedVariant = null
                        ),
                        Exercise.ChooseTranslation(
                            word = "Journey",
                            answerVariants = listOf("Путешествие", "Работа", "Отдых"),
                            correctVariant = 0,
                            selectedVariant = null
                        ),
                        Exercise.InputWord(
                            translation = "Улучшать",
                            correctAnswer = "improve",
                            inputtedWord = null
                        ),
                    )
                )
            }
        ),

        Lesson(
            lessonId = "lesson_real_02",
            components = buildList {
                // New Words
                addAll(
                    listOf(
                        Exercise.NewWord(
                            word = "Achieve",
                            translation = "Достигать",
                            phoneticTranscription = "/əˈtʃiːv/",
                            doesUserKnow = null,
                            examples = listOf(
                                "He worked hard to achieve his goals." to "Он усердно работал, чтобы достичь своих целей.",
                                "We can achieve success together." to "Мы можем достичь успеха вместе."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Support",
                            translation = "Поддержка",
                            phoneticTranscription = "/səˈpɔːt/",
                            doesUserKnow = null,
                            examples = listOf(
                                "Family support is important." to "Поддержка семьи важна.",
                                "Thank you for your support." to "Спасибо за вашу поддержку."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Opportunity",
                            translation = "Возможность",
                            phoneticTranscription = "/ˌɒp.əˈtʃuː.nə.ti/",
                            doesUserKnow = null,
                            examples = listOf(
                                "This is a great opportunity." to "Это отличная возможность.",
                                "Don't miss the opportunity." to "Не упусти возможность."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Challenge",
                            translation = "Вызов, трудность",
                            phoneticTranscription = "/ˈtʃæl.ɪndʒ/",
                            doesUserKnow = null,
                            examples = listOf(
                                "Learning a new language is a challenge." to "Изучение нового языка — это вызов.",
                                "He enjoys a good challenge." to "Ему нравятся хорошие испытания."
                            )
                        ),
                        Exercise.NewWord(
                            word = "Success",
                            translation = "Успех",
                            phoneticTranscription = "/səkˈses/",
                            doesUserKnow = null,
                            examples = listOf(
                                "Hard work leads to success." to "Усердная работа приводит к успеху.",
                                "Success comes from persistence." to "Успех приходит благодаря настойчивости."
                            )
                        )
                    )
                )

                // Exercises
                addAll(
                    listOf(
                        Exercise.ChooseTranslation(
                            word = "Achieve",
                            answerVariants = listOf("Игнорировать", "Достигать", "Забывать"),
                            correctVariant = 1,
                            selectedVariant = null
                        ),
                        Exercise.ChooseTranslation(
                            word = "Support",
                            answerVariants = listOf("Поддержка", "Конфликт", "Ошибка"),
                            correctVariant = 0,
                            selectedVariant = null
                        ),
                        Exercise.InputWord(
                            translation = "Вызов, трудность",
                            correctAnswer = "challenge",
                            inputtedWord = null
                        ),
                        Exercise.FillTheGap(
                            sentence = listOf("Hard work leads to ", "."),
                            answerVariants = listOf("success", "failure", "confusion"),
                            correctVariant = 0,
                            selectedVariant = null
                        )
                    )
                )

            }
        )
    )
}