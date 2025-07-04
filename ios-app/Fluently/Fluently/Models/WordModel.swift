//
//  WordModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

final class WordModel: Codable{
    var cefrLevel: String
    var exercise: ExerciseModel
    var isLearned: Bool
    var sentences: [SentenceModel]
    var subtopic: String
    var topic: String
    var transcription: String
    var translation: String
    var word: String
    var wordId: String

    init(
        cefrLevel: String,
        exercise: ExerciseModel,
        isLearned: Bool,
        sentences: [SentenceModel],
        subtopic: String,
        topic: String,
        transcription: String,
        translation: String,
        word: String,
        wordId: String
    ) {
        self.cefrLevel = cefrLevel
        self.exercise = exercise
        self.isLearned = isLearned
        self.sentences = sentences
        self.subtopic = subtopic
        self.topic = topic
        self.transcription = transcription
        self.translation = translation
        self.word = word
        self.wordId = wordId
    }

    enum CodingKeys: String, CodingKey {
        case cefrLevel = "cefr_level"
        case exercise
        case isLearned = "is_learned"
        case sentences
        case subtopic
        case topic
        case transcription
        case translation
        case word
        case wordId = "word_id"
    }
}
