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

}

final class DictionaryScreenPresenter: DictionaryScreenPresenting {
#if targetEnvironment(simulator)
    @Published var words: [Word] = Word.generateMockWords()
#else
    @Published var words: [Word] = Word.generateMockWords()
#endif
    @Published var filteredWords: [Word] = []

    init() {

        self.filteredWords = words
    }

    func filter(prefix: String) {
        guard !prefix.isEmpty else {
            filteredWords = words
            return
        }

        filteredWords = words.filter { $0.word.contains(prefix) }
    }
}
