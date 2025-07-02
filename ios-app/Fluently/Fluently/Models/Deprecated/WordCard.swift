//
//  WordCardExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

// check new word
final class WordCard: Exercise {
    // MARK: - Properties
    var wordId: UUID
    var word: String
    var translation: String
    var transcription: String
    var cefrLevel: String
    var isNew: Bool
    var topic: String
    var subtopic: String
    var sentences: [Sentence]
    var exercise: Exercise

    // MARK: - Init
    init(
        exerciseId: UUID,
        wordId: UUID,
        word: String,
        translation: String,
        transcription: String,
        cefrLevel: String,
        isNew:Bool,
        topic: String,
        subtopic: String,
        sentences: [Sentence],
        exercise: Exercise
    ) {
        self.wordId = wordId
        self.word = word
        self.translation = translation
        self.transcription = transcription
        self.cefrLevel = cefrLevel
        self.isNew = isNew
        self.topic = topic
        self.subtopic = subtopic
        self.sentences = sentences
        self.exercise = exercise

        super.init(
            exerciseId: exerciseId,
            exerciseType: "wordCard",
            correctAnswer: word
        )
    }
}

// sentence class
final class Sentence {
    // MARK: - Properties
    var sentenceId: UUID
    var sentence: String
    var translation: String

    // MARK: - Init
    init(
        sentenceId: UUID,
        sentece: String,
        translation: String
    ) {
        self.sentenceId = sentenceId
        self.sentence = sentece
        self.translation = translation
    }
}
