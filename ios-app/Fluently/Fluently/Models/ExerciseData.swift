//
//  ExerciseData.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.07.2025.
//

import Foundation
import SwiftData

//@Model
//class ExerciseData: Codable{
//    var correctAnswer: String
//
//    init(correctAnswer: String) {
//        self.correctAnswer = correctAnswer
//    }
//
//    // MARK: - Codable
//    private enum CodingKeys: String, CodingKey {
//        case correctAnswer = "correct_answer"
//    }
//
//    required init(from decoder: Decoder) throws {
//        let container = try decoder.container(keyedBy: CodingKeys.self)
//        correctAnswer = try container.decode(String.self, forKey: .correctAnswer)
//    }
//
//    func encode(to encoder: Encoder) throws {
//        var container = encoder.container(keyedBy: CodingKeys.self)
//        try container.encode(correctAnswer, forKey: .correctAnswer)
//    }
//}
//
