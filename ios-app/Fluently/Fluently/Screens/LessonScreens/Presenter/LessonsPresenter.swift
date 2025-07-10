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
    private(set) var words: [WordModel]
    @Published private(set) var currentExNumber: Int
    @Published private(set) var currentEx: ExerciseModel
    @Published private(set) var currentExType: ExerciseModelType

    @Published private(set) var learned = 0
    private(set) var wordsPerLesson = 10

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
        currentExType = currentEx.type
    }

    func answer(_ answer: String) {
        if currentEx.data.correctAnswer.lowercased() == answer.lowercased() {
            statistic[.correct]!.append(currentEx)
        } else {
            statistic[.uncorrect]!.append(currentEx)
        }
        nextExercise()
        learned += 1
    }

    // MARK: - Lesson navigation
    func nextExercise() {
        guard currentExNumber < words.count - 1 else {
            finishLesson()
            return
        }

        if learned == 9 {
            finishLesson()
        }
        currentExNumber += 1
        currentEx = words[currentExNumber].exercise
        currentExType = .wordCard
    }

    // func to represent statistic
    func finishLesson() {
        navigateBack()
        statistic.keys.forEach { solution in
            print("------------ \(solution.rawValue) ------------")
            statistic[solution]?.forEach { exr in
                if let pickoptions = exr.data as? PickOptionSentence {
                    print("Pick option sentence: \(pickoptions.template) -> \(pickoptions.correctAnswer)")
                }
                if let chooseTranslation = exr.data as? ChooseTranslationEngRuss {
                    print("Choose translation: \(chooseTranslation.text) -> \(chooseTranslation.correctAnswer)")
                }
                if let typeTranslation = exr.data as? WriteFromTranslation {
                    print("Type translation: \(typeTranslation.translation) -> \(typeTranslation.correctAnswer)")
                }
            }
        }
        print("--------------------------------------")
    }
}

