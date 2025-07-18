//
//  DictionaryScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI
import SwiftData

protocol DictionaryScreenPresenting: ObservableObject {

    func filter(prefix: String)
}

final class DictionaryScreenPresenter: DictionaryScreenPresenting {

    var isLearned = false
#if targetEnvironment(simulator)
    @Published var words: [WordModel] = WordModel.generateMockWords()
#else
    var words: [WordModel] {
        guard let modelContext else { return [] }
        let descriptor = FetchDescriptor<WordModel>(
            predicate: #Predicate {
                $0.isLearned == isLearned &&
                $0.isInLesson == false &&
                $0.isDayWord == false
            },
            sortBy: [SortDescriptor(\.wordDate, order: .reverse)]
        )
        return (try? modelContext.fetch(descriptor)) ?? []
    }
#endif
    @Published var filteredWords: [WordModel] = []

    var modelContext: ModelContext?

    init(modelContext: ModelContext? = nil, isLearned: Bool) {
        self.isLearned = isLearned
        self.modelContext = modelContext
        self.filteredWords = words
    }

    func setModelContext(_ context: ModelContext) {
        self.modelContext = context
//        words.forEach { word in
//            word.isInLesson = false
//
//            if word.wordId == nil || word.word == nil {
//                modelContext?.delete(word)
//            }
//        }
//
//        try? modelContext?.save()
    }

    func filter(prefix: String) {
        guard !prefix.isEmpty else {
            filteredWords = words
            return
        }

        filteredWords = words.filter { $0.word!.contains(prefix.lowercased()) || $0.translation!.contains(prefix.lowercased())}
    }
}
