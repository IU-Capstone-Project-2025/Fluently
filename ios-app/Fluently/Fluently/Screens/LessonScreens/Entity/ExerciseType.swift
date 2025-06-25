//
//  LessonModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation

enum ExerciseType: String, CaseIterable {
    case chooseTranslationEngRuss = "chooseTranslationEngRuss"
    case typeTranslationRussEng = "typeTranslationRussEng"
    case pickOptions = "pickOptions"
    case recordPronounce = "recordPronounce"
    case wordCard = "wordCard"

    case numberOfWords = "numberOfWords"
}

enum ExerciseSolution: String{
    case correct = "correct"
    case uncorrect = "uncorrect"
}


class Exercise {
    var exerciseId: UUID
    var exerciseType: ExerciseType
    var correctAnswer: String

    init(exerciseId: UUID, exerciseType: String, correctAnswer: String) {
        self.exerciseId = exerciseId
        guard let type = ExerciseType(rawValue: exerciseType) else {
            fatalError("Invalid LessonType string: \(exerciseType)")
        }
        self.exerciseType = type
        self.correctAnswer = correctAnswer
    }
}
