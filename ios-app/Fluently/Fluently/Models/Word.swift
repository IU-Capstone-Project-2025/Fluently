//
//  Word.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import Foundation
import SwiftData

// TODO: Maybe delete ? 
@Model
class Word {
    // MARK: - Properties
    var id: UUID
    var word: String
    var translation: String
    var wordClass: String
    var context: String
    var transcription: String
    var sentences: [String]
    var topic_id: UUID

    // MARK: - Init
    init(
        id: UUID,
        word: String,
        translation: String,
        wordClass: String,
        context: String,
        transcription: String,
        sentences: [String],
        topic_id: UUID
    ) {
        self.id = id
        self.word = word
        self.translation = translation
        self.wordClass = wordClass
        self.context = context
        self.transcription = transcription
        self.sentences = sentences
        self.topic_id = topic_id
    }

    // generator mock word
    static func mockWord() -> Word {
        return Word(
            id: UUID(),
            word: "Car",
            translation: "Машина",
            wordClass: "Vechile",
            context: "Bob is Driving his car",
            transcription: "kar",
            sentences: [
                "Bob is driving a car",
                "I will by car next week"
            ],
            topic_id: UUID()
        )
    }

    static func generateMockWords() -> [Word] {
        return [
            Word(
                id: UUID(),
                word: "Cat",
                translation: "Кошка",
                wordClass: "Animal",
                context: "Bob is driving his car",
                transcription: "kɑːt",
                sentences: [
                    "John love his cat",
                    "The cat is a predator"
                ],
                topic_id: UUID()
            ),
            Word(
                id: UUID(),
                word: "Car",
                translation: "Машина",
                wordClass: "Vehicle",
                context: "Bob is driving his car",
                transcription: "kɑːr",
                sentences: [
                    "Bob is driving a car",
                    "I will buy a car next week"
                ],
                topic_id: UUID()
            ),
            Word(
                id: UUID(),
                word: "Book",
                translation: "Книга",
                wordClass: "Object",
                context: "She's reading an interesting book",
                transcription: "bʊk",
                sentences: [
                    "This book changed my life",
                    "Could you recommend me a good book?"
                ],
                topic_id: UUID()
            ),
            Word(
                id: UUID(),
                word: "Run",
                translation: "Бегать",
                wordClass: "Verb",
                context: "I run every morning",
                transcription: "rʌn",
                sentences: [
                    "He runs faster than me",
                    "Let's run together on Sunday"
                ],
                topic_id: UUID()
            ),
            Word(
                id: UUID(),
                word: "Beautiful",
                translation: "Красивый",
                wordClass: "Adjective",
                context: "What a beautiful sunset!",
                transcription: "ˈbjuːtɪfəl",
                sentences: [
                    "She has a beautiful voice",
                    "This is the most beautiful place I've ever seen"
                ],
                topic_id: UUID()
            ),
            Word(
                id: UUID(),
                word: "Quickly",
                translation: "Быстро",
                wordClass: "Adverb",
                context: "He finished his work quickly",
                transcription: "ˈkwɪkli",
                sentences: [
                    "Time passes quickly when you're having fun",
                    "Please respond quickly to this email"
                ],
                topic_id: UUID()
            )
        ]
    }
}
