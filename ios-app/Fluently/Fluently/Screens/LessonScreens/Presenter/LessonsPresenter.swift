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
    @Published private(set) var words: [WordModel]
    @Published private(set) var currentWordNumber: Int = 0
    
    @Published private(set) var currentEx: ExerciseModel
    @Published private(set) var currentExType: ExerciseModelType

    @Published private(set) var learnedCount = 0
    private(set) var wordsPerLesson = 10

    var statistic: [ExerciseSolution: [ExerciseModel]] = [:]
    var wordsProgress: [ExerciseSolution: [WordModel]] = [:]

    private var lessonsStack: [WordModel] = []
    private(set) var currentExerciseNumber: Int = 0

    @Published private(set) var isAIChat = false

    // MARK: - Init
    init(router: AppRouter) {
        self.router = router
        self.words = []

        self.currentEx = ExerciseModel.init(data: EmptyExerciseData(), type: .wordCard)
        self.currentExType = .wordCard

        /// Initialize statistics dictionaries
        statistic[.correct] = []
        statistic[.uncorrect] = []
        wordsProgress[.correct] = []
        wordsProgress[.uncorrect] = []
    }

    func fetchWords() throws {
        guard let context = modelContext else {
            print("No model")
            fatalError()
        }

        let predicate = #Predicate<WordModel> { word in
            word.isInLesson == true
        }

        let descriptor = FetchDescriptor<WordModel>(
            predicate: predicate,
            sortBy: [SortDescriptor(\.wordDate)]
        )

        words = try context.fetch(descriptor)
        currentEx = words[0].exercise!
    }

    // MARK: - Navigation

    func navigateBack() {
        router.pop()
    }

    /// Show exercises from the lessons stack
    private func showExercises() {
        guard !lessonsStack.isEmpty else { return }

        currentEx = lessonsStack[currentExerciseNumber].exercise!
        currentExType = currentEx.type
    }

    /// Add word to learning stack and handle navigation
    func willLearn() {
        guard currentWordNumber < words.count else { return }

        lessonsStack.append(words[currentWordNumber])
        learnedCount += 1

        if learnedCount == wordsPerLesson {
            showExercises()
        }

        /// 1. When we've added 3 words (batch learning)
        /// 2. When we've added all words for the lesson (10 words)
        /// 3. When we've processed all words in the list
        if lessonsStack.count % 3 == 0 || lessonsStack.count == wordsPerLesson || currentWordNumber == words.count - 1 {
            showExercises()
        } else {
            /// Next word card if we're not showing exercises yet
            currentWordNumber += 1
            if currentWordNumber < words.count {
                if let ex = words[currentWordNumber].exercise {
                    currentEx = ex
                    currentExType = .wordCard
                } else {
                    print(words[currentWordNumber].word ?? "Nil name")
                    print(words[currentWordNumber].exercise?.type ?? "Nil lesson")
                    print(currentWordNumber)
                    alreadyKnow()
                }
            }
        }
    }

    func alreadyKnow() {
        guard currentWordNumber < words.count else { return }

        wordsProgress[.correct, default: []].append(words[currentWordNumber])
        currentWordNumber += 1

        if currentWordNumber < words.count {
            if let ex = words[currentWordNumber].exercise {
                currentEx = ex
                currentExType = .wordCard
            } else {
                alreadyKnow()
            }
        } else if !lessonsStack.isEmpty && currentExerciseNumber != lessonsStack.count {
            /// If we finished words but have exercises to complete
            showExercises()
        } else {
            /// No words left and no exercises - finish lesson
            finishLesson()
        }
    }

    func answer(_ answer: String) {
        print("\(currentExerciseNumber), \(lessonsStack[currentExerciseNumber].word!), answer: \(answer)")
        guard !lessonsStack.isEmpty, currentExerciseNumber < lessonsStack.count else {
            print("Skip")
            return
        }

        let exWord = lessonsStack[currentExerciseNumber]
        let isCorrect = currentEx.exerciseData.correctAnswer.lowercased() == answer.lowercased()

        exWord.isLearned = isCorrect
        let solution: ExerciseSolution = isCorrect ? .correct : .uncorrect
        statistic[solution, default: []].append(currentEx)
        wordsProgress[solution, default: []].append(exWord)

        nextExercise()
    }

    // MARK: - Lesson navigation
    func nextExercise() {
        guard !lessonsStack.isEmpty else {
            finishLesson()
            return
        }

        currentExerciseNumber += 1

        if learnedCount == wordsPerLesson {
            finishLesson()
        }

        if currentExerciseNumber >= lessonsStack.count {
            if currentWordNumber < words.count - 1 {
                currentWordNumber += 1
                if let ex = words[currentWordNumber].exercise {
                    currentEx = ex
                    currentExType = .wordCard
                } else {
                    nextExercise()
                }
            } else {
                finishLesson()
            }
        } else {
            // Show next exercise in current batch
            currentEx = lessonsStack[currentExerciseNumber].exercise!
            currentExType = currentEx.type
        }
    }

    func finishLesson() {
        // Save progress
        words.forEach { word in
            word.isInLesson = false
            modelContext?.insert(word)
        }
        try? modelContext?.save()

        // Send progress to server
        let api = APIService()
        Task {
            do {
                try await api.sendProgress(words: wordsProgress[.correct] ?? [])
            } catch {
                print("Error saving process: \(error.localizedDescription)")
            }
        }

        // Print statistics
        printLessonStatistics()

        isAIChat = true
    }

    func closeLesson() {
        // Navigate back
        navigateBack()
    }

    private func printLessonStatistics() {
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

