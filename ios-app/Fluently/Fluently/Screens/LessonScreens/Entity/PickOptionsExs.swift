//
//  PickOptionsExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftData

final class PickOptionsExs: Exercise {
    var sentenceId: UUID
    var sentence: String
    var options: [String]

    init(
        exerciseId: UUID,
        sentenceId: UUID,
        sentence: String,
        options: [String],
        correctAnswer: String
    ) {
        self.sentenceId = sentenceId
        self.sentence = sentence
        self.options = options

        super.init(
            exerciseId: exerciseId,
            exerciseType: "pickOptions",
            correctAnswer: correctAnswer
        )
    }
}

struct PickOptionsGenerator {
    static func generateMockPickOptionsLessons() -> [PickOptionsExs] {
        return [
            PickOptionsExs(
                exerciseId: UUID(),
                sentenceId: UUID(),
                sentence: "The chemical symbol for gold is ____ ",
                options: ["Au", "Ag", "Go", "Gd"],
                correctAnswer: "Au"
            ),
            PickOptionsExs(
                exerciseId: UUID(),
                sentenceId: UUID(),
                sentence: "The largest planet in our solar system is ____ ",
                options: ["Earth", "Saturn", "Jupiter", "Neptune"],
                correctAnswer: "Jupiter"
            ),
            PickOptionsExs(
                exerciseId: UUID(),
                sentenceId: UUID(),
                sentence: "The programming language developed by Apple is ____ ",
                options: ["Java", "Swift", "Kotlin", "Dart"],
                correctAnswer: "Swift"
            ),
            PickOptionsExs(
                exerciseId: UUID(),
                sentenceId: UUID(),
                sentence: "The capital of Japan is ____ ",
                options: ["Beijing", "Seoul", "Tokyo", "Bangkok"],
                correctAnswer: "Tokyo"
            ),
            PickOptionsExs(
                exerciseId: UUID(),
                sentenceId: UUID(),
                sentence: "The longest river in the world is ____ ",
                options: ["Amason", "Nile", "Yangtse", "Mississippi"],
                correctAnswer: "Nile"
            )
        ]
    }
}
