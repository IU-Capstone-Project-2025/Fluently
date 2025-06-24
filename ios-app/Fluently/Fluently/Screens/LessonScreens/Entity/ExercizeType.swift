//
//  LessonModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation

enum ExercizeType: String {
    case chooseTranslationEngRuss = "chooseTranslationEngRuss"
    case chooseTranslationRussEng = "chooseTranslationRussEng"
    case pickOptions = "pickOptions"
    case recordPronounce = "recordPronounce"
    case wordCard = "wordCard"

    case numberOfWords = "numberOfWords"
}

enum ExercizeSolution: String{
    case correct = "correct"
    case uncorrect = "uncorrect"
}


class Exercize {
    var exercizeID: UUID
    var exercizeType: ExercizeType
    var correctAnswer: String

    init(exercizeID: UUID, exercizeType: String, correctAnswer: String) {
        self.exercizeID = exercizeID
        guard let type = ExercizeType(rawValue: exercizeType) else {
            fatalError("Invalid LessonType string: \(exercizeType)")
        }
        self.exercizeType = type
        self.correctAnswer = correctAnswer
    }
}
