//
//  WordCardmockGenerator.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

extension WordModel {
    static func mockWord() -> WordModel {
        return WordModel(
//            cefrLevel: "A1",
            exercise: ExerciseModel(
                data: PickOptionSentence(
                    template: "Bob is driving a car",
                    options: [
                        "plane",
                        "dog",
                        "cat",
                    ],
                    correctAnswer: "car"),
                type: .pickOptionSentence
            ),
            isLearned: false,
            sentences: [],
            subtopic: "Car",
            topic: "Vechile",
            transcription: "ka:r",
            translation: "Машина",
            word: "Car",
            wordId: UUID().uuidString
        )
    }

    static func generateMockWords(count: Int = 5) -> [WordModel] {
        let mockWords = [
            WordModel(
//                cefrLevel: "A1",
                exercise: ExerciseModel(
                    data: PickOptionSentence(
                        template: "The __________ is red",
                        options: [
                            "banana",
                            "watermelon",
                            "orange",
                            "apple"
                        ],
                        correctAnswer: "apple"
                    ),
                    type: .pickOptionSentence
                ),
                isLearned: false,
                sentences: [
                    SentenceModel(
                        text: "I eat an apple every morning",
                        translation: "Я ем яблоко каждое утро"
                    )
                ],
                subtopic: "Food",
                topic: "Fruits",
                transcription: "ˈæp.əl",
                translation: "Яблоко",
                word: "Apple",
                wordId: UUID().uuidString
            ),
            WordModel(
//                cefrLevel: "A2",
                exercise: ExerciseModel(
                    data: ChooseTranslationEngRuss(
                        text: "Книга",
                        options: [
                            "magazine",
                            "newspaper",
                            "notebook",
                            "book"
                        ],
                        correctAnswer: "book"
                    ),
                    type: .typeTranslationRussEng
                ),
                isLearned: true,
                sentences: [
                    SentenceModel(
                        text: "The book is on the table",
                        translation: "Книга на столе"
                    ),
                    SentenceModel(
                        text: "I love reading books",
                        translation: "Я люблю читать книги"
                    )
                ],
                subtopic: "Education",
                topic: "Objects",
                transcription: "bʊk",
                translation: "Книга",
                word: "Book",
                wordId: UUID().uuidString
            ),
            WordModel(
//                cefrLevel: "B1",
                exercise: ExerciseModel(
                    data: WriteFromTranslation(
                        translation: "Бегать",
                        correctAnswer: "Run"
                    ),
                    type: .typeTranslationRussEng
                ),
                isLearned: false,
                sentences: [
                    SentenceModel(
                        text: "I run in the park every weekend",
                        translation: "Я бегаю в парке каждые выходные"
                    )
                ],
                subtopic: "Activities",
                topic: "Sports",
                transcription: "rʌn",
                translation: "Бегать",
                word: "Run",
                wordId: UUID().uuidString
            ),
            WordModel(
//                cefrLevel: "B2",
                exercise: ExerciseModel(
                    data: WriteFromTranslation(
                        translation: "Красивый",
                        correctAnswer: "Beautiful"
                    ),
                    type: .typeTranslationRussEng
                ),
                isLearned: true,
                sentences: [
                    SentenceModel(
                        text: "This is a beautiful painting",
                        translation: "Это прекрасная картина"
                    )
                ],
                subtopic: "Art",
                topic: "Adjectives",
                transcription: "ˈbjuː.tɪ.fəl",
                translation: "Красивый",
                word: "Beautiful",
                wordId: UUID().uuidString
            )
        ]

        if count <= mockWords.count {
            return Array(mockWords.prefix(count))
        } else {
            var words = mockWords
            for _ in mockWords.count..<count {
                words.append(mockWord())
            }
            return words
        }
    }
}
