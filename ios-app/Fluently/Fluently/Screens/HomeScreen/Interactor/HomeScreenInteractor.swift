//
//  HomeScreenInteractor.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.07.2025.
//

import Foundation


final class HomeScreenInteractor {

    let api: APIService = APIService()

    func getLesson() {
        Task {
            do {
                try await api.getLessons()
            } catch {
                print("Error: \(error)")
            }
        }
    }
}
