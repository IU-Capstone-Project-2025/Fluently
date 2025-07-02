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
        subTopic: String,
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
        self.subtopic = subTopic
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

extension WordModel {
    static func mockWord() -> WordModel {
        return WordModel(
            cefrLevel: "A1",
            exercise: ExerciseModel(
                data: "",
                type: ""
            ),
            isLearned: false,
            sentences: [],
            subTopic: "Car",
            topic: "Vechile",
            transcription: "ka:r",
            translation: "Машина",
            word: "Car",
            wordId: UUID().uuidString
        )
    }

    static func generateMockWords() -> [WordModel] {
        return []
    }
}


final class SentenceModel: Codable{
    var id: UUID
    var text: String
    var translation: String

    init(
        text: String,
        translation: String
    ) {
        self.id = UUID()
        self.text = text
        self.translation = translation
    }
}

extension SentenceModel: Hashable {
    func hash(into hasher: inout Hasher) {
        hasher.combine(id)
    }

    static func == (lhs: SentenceModel, rhs: SentenceModel) -> Bool {
        return lhs.text.compare(rhs.text).rawValue == 0
    }
}

final class ExerciseModel: Codable{
    var data: String
    var type: String

    init(
        data: String,
        type: String
    ) {
        self.data = data
        self.type = type
    }
}
