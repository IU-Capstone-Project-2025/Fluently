//
//  WordModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import SwiftData

@Model
final class WordModel: Codable{
    var exercise: ExerciseModel
    var isLearned: Bool = false
    var sentences: [SentenceModel]
    var subtopic: String
    var topic: String
    var transcription: String = "kar"
    var translation: String
    var word: String
    var wordId: String

    init(
        exercise: ExerciseModel,
        isLearned: Bool,
        sentences: [SentenceModel],
        subtopic: String,
        topic: String,
//        transcription: String,
        translation: String,
        word: String,
        wordId: String
    ) {
        self.exercise = exercise
        self.sentences = sentences
        self.subtopic = subtopic
        self.topic = topic
//        self.transcription = transcription
        self.translation = translation
        self.word = word
        self.wordId = wordId
    }

    enum CodingKeys: String, CodingKey {
        case exercise
        case isLearned = "is_learned"
        case sentences
        case subtopic
        case topic
//        case transcription
        case translation
        case word
        case wordId = "word_id"
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        exercise = try container.decode(ExerciseModel.self, forKey: .exercise)
        isLearned = false
        sentences = try container.decode([SentenceModel].self, forKey: .sentences)
        subtopic = try container.decode(String.self, forKey: .subtopic)
        topic = try container.decode(String.self, forKey: .topic)
//        transcription: String,
        translation = try container.decode(String.self, forKey: .translation)
        word = try container.decode(String.self, forKey: .word)
        wordId = try container.decode(String.self, forKey: .wordId)
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(exercise, forKey: .exercise)
        try container.encode(isLearned, forKey: .isLearned)
        try container.encode(sentences, forKey: .sentences)
        try container.encode(subtopic, forKey: .subtopic)
        try container.encode(topic, forKey: .topic)
//        transcription: String,
        try container.encode(translation, forKey: .translation)
        try container.encode(word, forKey: .word)
        try container.encode(wordId, forKey: .wordId)
    }
}
