//
//  PickOptionsExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftData

@Model
// exr to match word with sentence
final class PickOptionSentence: ExerciseData {
    // MARK: - Properties
    var template: String
    var options: [String]
    var correctAnswer: String

    // MARK: - Init
    init(
        template: String,
        options: [String],
        correctAnswer: String
    ) {
        self.template = template
        self.options = options

        self.correctAnswer = correctAnswer
    }
    
    // MARK: - Codable
    private enum CodingKeys: String, CodingKey {
        case template
        case options = "pick_options"
        case correctAnswer = "correct_answer"
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        template = try container.decode(String.self, forKey: .template)
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
        try container.encode(template, forKey: .template)
        try container.encode(options, forKey: .options)
        try container.encode(correctAnswer, forKey: .correctAnswer)
    }
}
