//
//  HomeScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI
import SwiftData

// MARK: - Protocol for presenter
protocol HomeScreenPresenting: ObservableObject {
    func getLesson() async throws

    // Navigation
    func navigatoToProfile()
    func navigatoToLesson()
}

// MARK: - Presenter implementation
final class HomeScreenPresenter: HomeScreenPresenting {
    let router: HomeScreenRouter
    let interactor: HomeScreenInteractor

    @ObservedObject var account: AccountData
#if targetEnvironment(simulator)
    @Published var lesson: CardsModel? = CardsModel(
        cards: WordModel.generateMockWords(count: 5),
        lesson: LessonModel(
            startedAt: "",
            totalWords: 10,
            wordsPerLesson: 5,
            cefrLevel: "A1"
        )
    )

#else
    @Published var lesson: CardsModel?
#endif

    @Published var wordOfTheDay: WordModel?

    var modelContext: ModelContext?

    init(
        router: HomeScreenRouter,
        interactor: HomeScreenInteractor,
        account: AccountData
    ) {
        self.router = router
        self.interactor = interactor
        self.account = account
    }

    func getDayWord() {
        wordOfTheDay = getTodaysWord()
        guard wordOfTheDay == nil else {
            return
        }

        Task {
            do {
                self.wordOfTheDay = try await interactor.getDayWord()
                await saveWordOfTheDay()
            } catch {
                print("Error on getting word of the day: \(error.localizedDescription)")
                wordOfTheDay = WordModel.mockWord()
            }
        }
    }

    func saveWordOfTheDay() async {
        guard let modelContext else {
            return
        }
        let dayWordDTO = DayWord(
            word: wordOfTheDay
        )
        modelContext.insert(dayWordDTO)
    }

    func getTodaysWord() -> WordModel? {
        guard let modelContext else {
            print("no context")
            return WordModel.mockWord()
        }

        let today = Calendar.current.startOfDay(for: Date())
        let predicate = #Predicate<DayWord> {
            $0.date >= today
        }

        let descriptor = FetchDescriptor<DayWord>(predicate: predicate)
        return try? modelContext.fetch(descriptor).first?.word
    }

    @MainActor
    func getLesson() async throws {
        guard lesson == nil else {
            return
        }

        lesson = try await interactor.getLesson()
    }

    // Builders 
    func buildNotesScreen() -> NotesView{
        return NotesScreenBuilder.build(router: router.router)
    }

    func buildDictionaryScreen(isLearned: Bool) -> DictionaryView{
        return DictionaryScreenBuilder.build(
            isLearned: isLearned
        )
    }

    // Navigation

    func navigatoToProfile() {
        router.navigatoToProfile()
    }

    func navigatoToLesson() {
        router.navigatoToLesson(lesson!)
        lesson = nil
    }
}
