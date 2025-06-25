//
//  TypeTranslationExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation


final class TypeTranslationExs: Exercise {
    var wordId: UUID
    var word: String

    init(
        exerciseId: UUID,
        wordId: UUID,
        word: String,
        correctAnswer: String
    ) {
        self.wordId = wordId
        self.word = word

        super.init(
            exerciseId: exerciseId,
            exerciseType: "typeTranslationRussEng",
            correctAnswer: correctAnswer
        )
    }
}
