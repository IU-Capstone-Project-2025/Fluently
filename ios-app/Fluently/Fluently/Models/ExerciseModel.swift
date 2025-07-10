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

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        type = try container.decode(ExerciseModelType.self, forKey: .type)

        switch type {
            case .chooseTranslationEngRuss:
                data = try ChooseTranslationEngRuss(from: container.superDecoder(forKey: .data))
            case .typeTranslationRussEng:
                data = try WriteFromTranslation(from: container.superDecoder(forKey: .data))
            case .pickOptionSentence:
                data = try PickOptionSentence(from: container.superDecoder(forKey: .data))
            case .recordPronounce:
                // Assuming you have a RecordPronounce class
    //            data = try RecordPronounce(from: container.superDecoder(forKey: .data))
                data = try PickOptionSentence(from: container.superDecoder(forKey: .data))
            case .wordCard:
    //            data = try WordCard(from: container.superDecoder(forKey: .data))
                data = try PickOptionSentence(from: container.superDecoder(forKey: .data))
            case .numberOfWords:
    //            data = try NumberOfWords(from: container.superDecoder(forKey: .data))
                data = try PickOptionSentence(from: container.superDecoder(forKey: .data))
        }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(type, forKey: .type)

        switch type {
        case .chooseTranslationEngRuss:
            try (data as! ChooseTranslationEngRuss).encode(to: container.superEncoder(forKey: .data))
        case .typeTranslationRussEng:
            try (data as! WriteFromTranslation).encode(to: container.superEncoder(forKey: .data))
        case .pickOptionSentence:
            try (data as! PickOptionSentence).encode(to: container.superEncoder(forKey: .data))
        case .recordPronounce:
                print("nothing")
//            try (data as! RecordPronounce).encode(to: container.superEncoder(forKey: .data))
        case .wordCard:
                print("nothing")
//            try (data as! WordCard).encode(to: container.superEncoder(forKey: .data))
        case .numberOfWords:
                print("nothing")
//            try (data as! NumberOfWords).encode(to: container.superEncoder(forKey: .data))
        }
    }

    enum CodingKeys: String, CodingKey {
        case data
        case type
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
