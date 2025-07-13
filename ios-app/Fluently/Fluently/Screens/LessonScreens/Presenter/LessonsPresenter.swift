//
//  LessonsPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI
import SwiftData

final class LessonsPresenter: ObservableObject {
    // MARK: - Key Object
    private var router: AppRouter

    var modelContext: ModelContext?

    // MARK: - Properties
    private(set) var words: [WordModel]
    @Published private(set) var currentExNumber: Int
    @Published private(set) var currentEx: ExerciseModel
    @Published private(set) var currentExType: ExerciseModelType

    @Published private(set) var learned = 0
    private(set) var wordsPerLesson = 10

    var statistic: [ExerciseSolution : [ExerciseModel]]
    var wordsProgress: [ExerciseSolution: [WordModel]] = [:]

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

        wordsProgress[.correct] = []
        wordsProgress[.uncorrect] = []
    }

    // MARK: - Navigation

    func navigateBack() {
        router.pop()
    }

    func showLesson() {
        currentEx = words[currentExNumber].exercise
        currentExType = currentEx.type
    }

    func alreadyKnow() {
        currentExNumber += 1
        currentEx = words[currentExNumber].exercise
        wordsProgress[.correct]!.append(words[currentExNumber])
        currentExType = .wordCard
    }

    func answer(_ answer: String) {
        if currentEx.exerciseData.correctAnswer.lowercased() == answer.lowercased() {
            words[currentExNumber].isLearned = true
            statistic[.correct]!.append(currentEx)
            wordsProgress[.correct]!.append(words[currentExNumber])
        } else {
            words[currentExNumber].isLearned = false
            statistic[.uncorrect]!.append(currentEx)
            wordsProgress[.uncorrect]!.append(words[currentExNumber])
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
        words.forEach { word in
            modelContext?.insert(word)
        }
        try? modelContext?.save()

        let api = APIService()

        Task {
            try? await api.sendProgress(words: wordsProgress[.correct]!)
        }

        navigateBack()
        statistic.keys.forEach { solution in
            print("------------ \(solution.rawValue) ------------")
            statistic[solution]?.forEach { exr in
                if let pickoptions = exr.exerciseData as? PickOptionSentence {
                    print("Pick option sentence: \(pickoptions.template) -> \(pickoptions.correctAnswer)")
                }
                if let chooseTranslation = exr.exerciseData as? ChooseTranslationEngRuss {
                    print("Choose translation: \(chooseTranslation.text) -> \(chooseTranslation.correctAnswer)")
                }
                if let typeTranslation = exr.exerciseData as? WriteFromTranslation {
                    print("Type translation: \(typeTranslation.translation) -> \(typeTranslation.correctAnswer)")
                }
            }
        }
        print("--------------------------------------")
    }
}

