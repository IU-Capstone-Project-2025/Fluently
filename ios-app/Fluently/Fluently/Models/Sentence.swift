//
//  Sentence.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

// sentence class
final class Sentence {
    // MARK: - Properties
    var sentenceId: UUID
    var sentence: String
    var translation: String

    // MARK: - Init
    init(
        sentenceId: UUID,
        sentece: String,
        translation: String
    ) {
        self.sentenceId = sentenceId
        self.sentence = sentece
        self.translation = translation
    }
}
