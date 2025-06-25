//
//  ChooseTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation

// exr to choose correct translation
final class ChooseTranslationExs: Exercise {
    // MARK: - Properties
    var wordId: UUID
    var word: String
    var options: [String]

    // MARK: - Init
    init(
        exerciseId: UUID,
        wordId: UUID,
        word: String,
        options: [String],
        correctAnswer: String
    ) {
        self.wordId = wordId
        self.word = word
        self.options = options

        super.init(
            exerciseId: exerciseId,
            exerciseType: "chooseTranslationEngRuss",
            correctAnswer: correctAnswer
        )
    }
}
