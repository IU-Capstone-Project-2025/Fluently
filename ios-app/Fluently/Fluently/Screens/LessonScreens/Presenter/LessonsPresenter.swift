//
//  LessonsPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

final class LessonsPresenter: ObservableObject {

    private var router: AppRouter

    private var exercizes: [Exercize]
    @Published private(set) var currentExNumber: Int
    @Published private(set) var currentEx: Exercize

    var statistic: [ExercizeSolution : [Exercize]]

    init(router: AppRouter, exercizes: [Exercize]) {
        self.router = router

        self.exercizes = exercizes

        self.currentExNumber = 0
        self.currentEx = exercizes[0]
        self.statistic = [:]

        statistic[.correct] = []
        statistic[.uncorrect] = []
    }

    // MARK: - Navigation

    func navigateBack() {
        router.pop()
    }

    func answer(_ answer: String) {
        if currentEx.correctAnswer == answer {
            statistic[.correct]!.append(currentEx)
        } else {
            statistic[.uncorrect]!.append(currentEx)
        }
        nextExercize()
    }

    // MARK: - Lesson navigation
    func nextExercize() {
        guard currentExNumber < exercizes.count - 1 else {
            finishLesson()
            return
        }

        currentExNumber += 1
        currentEx = exercizes[currentExNumber]
    }

    func finishLesson() {
        navigateBack()
        statistic.keys.forEach { solution in
            print("------------ \(solution.rawValue) ------------")
            statistic[solution]?.forEach { exr in
                if let pickoptions = exr as? PickOptionsExs {
                    print(pickoptions.sentence)
                }
            }
        }
        print("--------------------------------------")
    }
}

