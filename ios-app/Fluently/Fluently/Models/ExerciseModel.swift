//
//  ExerciseModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import SwiftData

@Model
final class ExerciseModel: Codable {
    private var storedData: Data?
    var type: ExerciseModelType

    var exerciseData: any ExerciseData{
        get {
            guard let storedData else {
                return EmptyExerciseData()
            }
            let decoder = JSONDecoder()
            do {
                switch type {
                    case .chooseTranslationEngRuss:
                        return try decoder.decode(ChooseTranslationEngRuss.self, from: storedData)
                    case .typeTranslationRussEng:
                        return try decoder.decode(WriteFromTranslation.self, from: storedData)
                    case .pickOptionSentence:
                        return try decoder.decode(PickOptionSentence.self, from: storedData)
                    case .recordPronounce, .wordCard, .numberOfWords:
                        return try decoder.decode(PickOptionSentence.self, from: storedData)
                }
            } catch {
                fatalError("damn")
                return EmptyExerciseData()
            }
        }
        set {
            let encoder = JSONEncoder()
            storedData = try? encoder.encode(newValue)
        }
    }

    init(data: any ExerciseData, type: ExerciseModelType) {
        self.type = type
        self.exerciseData = data
    }

    private enum CodingKeys: String, CodingKey {
        case data
        case type
    }

    required init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        type = try container.decode(ExerciseModelType.self, forKey: .type)

        switch type {
            case .pickOptionSentence:
                let pickOptionData = try container.decode(PickOptionSentence.self, forKey: .data)
                self.exerciseData = pickOptionData
            case .chooseTranslationEngRuss:
                let translationData = try container.decode(ChooseTranslationEngRuss.self, forKey: .data)
                self.exerciseData = translationData
            case .typeTranslationRussEng:
                let writeData = try container.decode(WriteFromTranslation.self, forKey: .data)
                self.exerciseData = writeData
            default:
                let defaultData = try container.decode(PickOptionSentence.self, forKey: .data)
                self.exerciseData = defaultData
        }
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        try container.encode(type, forKey: .type)
        try container.encode(storedData, forKey: .data)
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
