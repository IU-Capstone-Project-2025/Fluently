//
//  CardsModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import SwiftData

// MARK: - Cards Model
/// model of data that returning at lesson request
@Model
final class CardsModel: Codable, Sendable {
    var cards: [WordModel]
    var lesson: LessonModel

    init(
        cards: [WordModel],
        lesson: LessonModel
    ) {
        self.cards = cards
        self.lesson = lesson
    }

    // MARK: - Codable
    enum CodingKeys: String, CodingKey {
        case cards
        case lesson
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        cards = try container.decode([WordModel].self, forKey: .cards)
        lesson = try container.decode(LessonModel.self, forKey: .lesson)
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(cards, forKey: .cards)
        try container.encode(lesson, forKey: .lesson)
    }
}
