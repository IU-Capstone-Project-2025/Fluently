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

    var statistic: [ExerciseSolution : [ExerciseModel]]

    // MARK: - Init
    init(router: AppRouter, words: [WordModel]) {
        self.router = router

        self.words = words

        self.currentExNumber = 0
        self.currentEx = words[0].exercise
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
                if let pickoptions = exr as? PickOptionsExs {
                    print(pickoptions.sentence)
                }
                if let chooseTranslation = exr as? ChooseTranslationExs {
                    print("Choose translation: \(chooseTranslation.word) -> \(chooseTranslation.correctAnswer)")
                }
                if let typeTranslation = exr as? TypeTranslationExs {
                    print("Type translation: \(typeTranslation.word) -> \(typeTranslation.correctAnswer)")
                }
            }
        }
        print("--------------------------------------")
    }
}

