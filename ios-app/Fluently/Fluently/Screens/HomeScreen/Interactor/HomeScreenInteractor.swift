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
        cards.cards.forEach { card in
            card.isInLibrary = false
        }
        return cards
    }

    func printCards(_ cards: CardsModel) {
        cards.cards.forEach { card in
            print(card.word ?? "Nil word")
            print(card.exercise?.type ?? "Nil exercise")
            print(type(of: card.exercise?.exerciseData))
        }
    }

    func getPrefs() async throws -> PreferencesModel {
        return try await api.getPreferences()
    }

    func getDayWord() async throws -> WordModel {
        return try await api.getDayWord()
    }
}
