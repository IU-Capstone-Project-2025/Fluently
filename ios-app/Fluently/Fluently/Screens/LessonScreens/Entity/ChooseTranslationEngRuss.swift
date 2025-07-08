//
//  ChooseTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

// exr to choose correct translation
final class ChooseTranslationEngRuss: ExerciseData {
    // MARK: - Properties
    var text: String
    var options: [String]

    // MARK: - Init
    init(text: String, options: [String], correctAnswer: String) {
        self.text = text
        self.options = options

        super.init(correctAnswer: correctAnswer)
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
        let temp = try container.decode(String.self, forKey: .correctAnswer)

//        options = try container.decode([String].self, forKey: .options)
        options = [temp]
        try super.init(from: container.superDecoder())
    }

    override func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(text, forKey: .text)
        try container.encode(options, forKey: .options)
        try super.encode(to: encoder)
    }
}
