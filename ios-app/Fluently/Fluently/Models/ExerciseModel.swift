//
//  ExerciseModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

class ExerciseModel: Codable{
    var data: ExerciseData
    var type: ExerciseModelType

    init(
        data: ExerciseData,
        type: String
    ) {
        self.data = data
        self.type = ExerciseModelType(rawValue: type) ?? .wordCard
    }
}

// MARK: - Exrs Types
enum ExerciseModelType: String, CaseIterable{
    case chooseTranslationEngRuss = "translate_ru_to_en"
    case typeTranslationRussEng = "write_word_from_translation"
    case pickOptionSentence = "pick_option_sentence"
    case recordPronounce = "recordPronounce"

    case wordCard = "word_card"
    case numberOfWords = "numberOfWords"
}

extension ExerciseModelType: Codable{}
