//
//  TypeTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftData

@Model
// exr: type correct translation
final class WriteFromTranslation: ExerciseData {
    // MARK: - Properties
    var translation: String
    var correctAnswer: String

    // MARK: - Init

    init(
        translation: String,
        correctAnswer: String
    ) {
        self.translation = translation

        self.correctAnswer = correctAnswer
    }
    
    // MARK: - Codable
    private enum CodingKeys: String, CodingKey {
        case translation
        case correctAnswer = "correct_answer"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        translation = try container.decode(String.self, forKey: .translation)
        correctAnswer = try container.decode(String.self, forKey: .correctAnswer)
//        try super.init(from: container.superDecoder())
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(translation, forKey: .translation)
        try container.encode(correctAnswer, forKey: .correctAnswer)
    }
}
