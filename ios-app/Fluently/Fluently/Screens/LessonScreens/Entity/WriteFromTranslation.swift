//
//  TypeTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

// exr: type correct translation
final class WriteFromTranslation: ExerciseData {
    // MARK: - Properties
    var translation: String

    // MARK: - Init

    init(
        translation: String,
        correctAnswer: String
    ) {
        self.translation = translation

        super.init(correctAnswer: correctAnswer)
    }
    
    // MARK: - Codable
    private enum CodingKeys: String, CodingKey {
        case translation
        case correctAnswer = "correct_answer"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        translation = try container.decode(String.self, forKey: .translation)
        try super.init(from: container.superDecoder())
    }

    override func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(translation, forKey: .translation)
        try super.encode(to: encoder)
    }
}
