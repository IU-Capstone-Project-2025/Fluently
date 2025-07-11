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
//        options = try container.decodeIfPresent([String].self, forKey: .options) ?? []
//        try super.init(from: container.superDecoder())

        options = []

        let answer = try container.decode(String.self, forKey: .correctAnswer)
        super.init(correctAnswer: answer)

        if let options = try container.decodeIfPresent([String].self, forKey: .options) {
            self.options = options.isEmpty ? [correctAnswer] : options
        } else {
            self.options = [correctAnswer]
        }
    }
    
    override func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(text, forKey: .text)
        try container.encode(options, forKey: .options)
        try super.encode(to: encoder)
    }
}
