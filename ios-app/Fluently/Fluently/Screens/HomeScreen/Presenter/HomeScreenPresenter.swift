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
    // Navigation
    func navigatoToProfile()
    func navigatoToLesson()

    // checkers for quality saved data
    func checkForNilIDs() async
    func compare()

    // daily word logic
    func saveWordOfTheDay() async
    func getTodaysWord() -> WordModel?

    // lesson logic
    func findLesson(context: ModelContext?) -> CardsModel?
    func getLesson() async throws
    func deleteLesson()

    // screens for cards 
    func buildDictionaryScreen(isLearned: Bool) -> DictionaryView
}

// MARK: - Presenter implementation
final class HomeScreenPresenter: HomeScreenPresenting {
    let router: HomeScreenRouter
    let interactor: HomeScreenInteractor

    @ObservedObject var account: AccountData
#if targetEnvironment(simulator)
    @Published var lesson: CardsModel? = CardsModel(
        cards: WordModel.generateMockWords(count: 20),
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
                if let wordOfTheDay {
                    wordOfTheDay.isDayWord = true
                    wordOfTheDay.isInLibrary = false
                }
                await saveWordOfTheDay()
            } catch {
                print("Error on getting word of the day: \(error.localizedDescription)")
                wordOfTheDay = WordModel.mockWord()
            }
        }
    }

    func checkForNilIDs() async {
        guard let modelContext else {
            return
        }

        let descriptor = FetchDescriptor<WordModel>()
        let words = try? modelContext.fetch(descriptor)

        guard let words else { return }
        var count = 0
        words.forEach { word in
            if word.wordId == nil {
                modelContext.delete(word)
                count += 1
            }
        }
        print("\(count) nil words deleted")
        try? modelContext.save()
    }

    func compare() {
        guard let modelContext else {
            return
        }

        guard let lesson else {
            return
        }

        let descriptor = FetchDescriptor<PreferencesModel>()
        let localPreferences = try? modelContext.fetch(descriptor).first

        var preferences: PreferencesModel?

        Task {
            if localPreferences == nil {
                preferences = try await interactor.getPrefs()
            } else {
                preferences = localPreferences
            }

            guard let preferences else {
                print("No valid preferences available")
                return
            }

            if lesson.lesson.wordsPerLesson < preferences.wordPerDay{
                print(lesson.lesson.wordsPerLesson,  preferences.wordPerDay)
                deleteLesson()
                try? await getLesson()
            } else {
                print("everthing ok")
            }
        }
    }

    func saveWordOfTheDay() async {
        guard let modelContext else {
            return
        }
        wordOfTheDay?.isDayWord = true

        let dayWordDTO = DayWord(
            word: wordOfTheDay
        )
        modelContext.insert(dayWordDTO)

        try? modelContext.save()
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
        if let existingLesson = findLesson(context: modelContext) {
            self.lesson = existingLesson
            return
        }

        let newLesson = try await interactor.getLesson()

        guard let modelContext else {
            throw LessonError.noModelContext
        }

        modelContext.insert(newLesson)
        try modelContext.save()

        self.lesson = newLesson

        print("Lesson saved in memory")
    }

    func findLesson(context: ModelContext?) -> CardsModel? {
        let descriptor = FetchDescriptor<CardsModel>()

        do {
            let found = try context?.fetch(descriptor).first
            if let found, found.cards.count == 0 { return nil }
            return found
        } catch {
            print("SwiftData fetch failed: \(error)")
            return nil
        }
    }

    func deleteLesson() {
        let descriptor = FetchDescriptor<CardsModel>()
        
        do {
            let lessons = try modelContext?.fetch(descriptor)

            if let lessons {
                lessons.forEach { l in
                    modelContext?.delete(l)
                }
            }
        } catch {
            print("SwiftData fetch failed: \(error)")
            return
        }

        if lesson != nil {
            modelContext?.delete(lesson!)
            lesson = nil
            print("lesson deleted")
        }
    }

    // Builders 
    func buildNotesScreen() -> some View{
#if targetEnvironment(simulator)
        return AIChatBuilder.build() {
            print("bebeb")
        }
#else
        return NotesScreenBuilder.build(router: router.router)
#endif
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
        deleteLesson()
    }


    enum LessonError: Error {
        case noModelContext
        case invalidLessonData
        case saveFailed(Error)
    }
}
