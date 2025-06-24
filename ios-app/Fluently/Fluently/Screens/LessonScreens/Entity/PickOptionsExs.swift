//
//  PickOptionsExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftData

final class PickOptionsExs: Exercize {
    @Attribute(.unique) var sentenceID: UUID
    var sentence: String
    var options: [String]

    init(exercizeID: UUID, exercizeType: String, sentenceID: UUID, sentence: String, options: [String], correctAnswer: String) {
        self.sentenceID = sentenceID
        self.sentence = sentence
        self.options = options

        super.init(
            exercizeID: exercizeID,
            exercizeType: exercizeType,
            correctAnswer: correctAnswer
        )
    }
}


struct PickOptionsGenerator {
    static func generateMockPickOptionsLessons() -> [PickOptionsExs] {
        return [
            PickOptionsExs(
                exercizeID: UUID(),
                exercizeType: "pickOptions",
                sentenceID: UUID(),
                sentence: "The chemical symbol for gold is ____ ",
                options: ["Au", "Ag", "Go", "Gd"],
                correctAnswer: "Au"
            ),
            PickOptionsExs(
                exercizeID: UUID(),
                exercizeType: "pickOptions",
                sentenceID: UUID(),
                sentence: "The largest planet in our solar system is ____ ",
                options: ["Earth", "Saturn", "Jupiter", "Neptune"],
                correctAnswer: "Jupiter"
            ),
            PickOptionsExs(
                exercizeID: UUID(),
                exercizeType: "pickOptions",
                sentenceID: UUID(),
                sentence: "The programming language developed by Apple is ____ ",
                options: ["Java", "Swift", "Kotlin", "Dart"],
                correctAnswer: "Swift"
            ),
            PickOptionsExs(
                exercizeID: UUID(),
                exercizeType: "pickOptions",
                sentenceID: UUID(),
                sentence: "The capital of Japan is ____ ",
                options: ["Beijing", "Seoul", "Tokyo", "Bangkok"],
                correctAnswer: "Tokyo"
            ),
            PickOptionsExs(
                exercizeID: UUID(),
                exercizeType: "pickOptions",
                sentenceID: UUID(),
                sentence: "The longest river in the world is ____ ",
                options: ["Amazon", "Nile", "Yangtze", "Mississippi"],
                correctAnswer: "Nile"
            )
        ]
    }
}
