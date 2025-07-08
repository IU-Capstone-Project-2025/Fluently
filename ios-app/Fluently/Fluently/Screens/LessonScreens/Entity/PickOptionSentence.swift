//
//  PickOptionsExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftData

// exr to match word with sentence
final class PickOptionSentence: ExerciseData {
    // MARK: - Properties
    var template: String
    var options: [String]

    // MARK: - Init
    init(
        template: String,
        options: [String],
        correctAnswer: String
    ) {
        self.template = template
        self.options = options

        super.init(correctAnswer: correctAnswer)
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
        let temp = try container.decode(String.self, forKey: .correctAnswer)

//        options = try container.decode([String].self, forKey: .options)
        options = [temp]
        try super.init(from: container.superDecoder())
    }

    override func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(template, forKey: .template)
        try container.encode(options, forKey: .options)
        try super.encode(to: encoder)
    }
}
