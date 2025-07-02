//
//  CardsModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

final class CardsModel: Codable {
    var cards: [WordModel]
    var lesson: LessonModel

    init(
        cards: [WordModel],
        lesson: LessonModel
    ) {
        self.cards = cards
        self.lesson = lesson
    }
}
