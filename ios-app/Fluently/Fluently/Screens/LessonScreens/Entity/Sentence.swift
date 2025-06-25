//
//  Sentence.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

final class Sentence {
    var sentenceId: UUID
    var sentence: String
    var translation: String

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
