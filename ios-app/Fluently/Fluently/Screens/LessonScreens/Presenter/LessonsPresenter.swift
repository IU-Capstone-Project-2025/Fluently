//
//  LessonsPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

final class LessonsPresenter: ObservableObject {
    // MARK: - Key Object
    private var router: AppRouter

    // MARK: - Properties
    private var words: [WordModel]
    @Published private(set) var currentExNumber: Int
    @Published private(set) var currentEx: ExerciseModel
    @Published private(set) var currentExType: ExerciseModelType

    var statistic: [ExerciseSolution : [ExerciseModel]]

    // MARK: - Init
    init(router: AppRouter, words: [WordModel]) {
        self.router = router

        self.words = words

        self.currentExNumber = 0
        self.currentEx = words[0].exercise
        self.currentExType = .wordCard
        self.statistic = [:]

        statistic[.correct] = []
        statistic[.uncorrect] = []
    }

    // MARK: - Navigation

    func navigateBack() {
        router.pop()
    }

    func showLesson() {
        currentEx = words[currentExNumber].exercise
        currentExType = ExerciseModelType(rawValue: currentEx.type) ?? .wordCard
    }

    func answer(_ answer: String) {
        if currentEx.correctAnswer == answer {
            statistic[.correct]!.append(currentEx)
        } else {
            statistic[.uncorrect]!.append(currentEx)
        }
        nextExercise()
    }

    // MARK: - Lesson navigation
    func nextExercise() {
        guard currentExNumber < words.count - 1 else {
            finishLesson()
            return
        }

        currentExNumber += 1
        currentEx = words[currentExNumber].exercise
    }

    // func to represent statistic
    func finishLesson() {
        navigateBack()
        statistic.keys.forEach { solution in
            print("------------ \(solution.rawValue) ------------")
            statistic[solution]?.forEach { exr in
                if let pickoptions = exr as? PickOptionSentence {
                    print(pickoptions.sentence)
                }
                if let chooseTranslation = exr as? chooseTranslationEngRuss {
                    print("Choose translation: \(chooseTranslation.word) -> \(chooseTranslation.correctAnswer)")
                }
                if let typeTranslation = exr as? WriteFromTranslation {
                    print("Type translation: \(typeTranslation.word) -> \(typeTranslation.correctAnswer)")
                }
            }
        }
        print("--------------------------------------")
    }
}

