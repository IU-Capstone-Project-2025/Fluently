//
//  Word.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import Foundation
import SwiftData

@Model
class Word {
    var id: UUID
    var word: String
    var translation: String
    var wordClass: String
    var context: String
    var transcription: String
    var topic_id: UUID

    init(id: UUID, word: String, translation: String, wordClass: String, context: String, transcription: String, topic_id: UUID) {
        self.id = id
        self.word = word
        self.translation = translation
        self.wordClass = wordClass
        self.context = context
        self.transcription = transcription
        self.topic_id = topic_id
    }

    static func mockWord() -> Word {
        return Word(
            id: UUID(),
            word: "Car",
            translation: "Машина",
            wordClass: "Vechile",
            context: "Bob is Driving his car",
            transcription: "kar",
            topic_id: UUID()
        )
    }
}
