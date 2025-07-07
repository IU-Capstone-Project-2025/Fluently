//
//  WordModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

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
        case sentences
        case subtopic
        case topic
//        case transcription
        case translation
        case word
        case wordId = "word_id"
    }
}
