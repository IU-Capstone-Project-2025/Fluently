//
//  HomeScreenInteractor.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.07.2025.
//

import Foundation


final class HomeScreenInteractor {

    let api: APIService = APIService()

    func getLesson() async throws -> CardsModel{
        let cards = try await api.getLesson()
        printCards(cards)
        return cards
    }

    func printCards(_ cards: CardsModel) {
        cards.cards.forEach { card in
            print(card.word)
            print(card.exercise.type)
            print(type(of: card.exercise.data))
        }
    }
}
