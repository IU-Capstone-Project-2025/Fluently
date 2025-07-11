//
//  ChooseTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftData

@Model
// exr to choose correct translation
final class ChooseTranslationEngRuss: ExerciseData {
    // MARK: - Properties
    var text: String
    var options: [String]
    var correctAnswer: String

    // MARK: - Init
    init(text: String, options: [String], correctAnswer: String) {
        self.text = text
        self.options = options
        self.correctAnswer = correctAnswer
    }
    
    // MARK: - Codable
    private enum CodingKeys: String, CodingKey {
        case text
        case options = "pick_options"
        case correctAnswer = "correct_answer"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        text = try container.decode(String.self, forKey: .text)
//        options = try container.decodeIfPresent([String].self, forKey: .options) ?? []
//        try super.init(from: container.superDecoder())

        options = []

        correctAnswer = try container.decode(String.self, forKey: .correctAnswer)


        if let options = try container.decodeIfPresent([String].self, forKey: .options) {
            self.options = options.isEmpty ? [correctAnswer] : options
        } else {
            self.options = [correctAnswer]
        }
    }
    
    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(text, forKey: .text)
        try container.encode(options, forKey: .options)
        try container.encode(correctAnswer, forKey: .correctAnswer)
    }
}
