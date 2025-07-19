//
//  ProfileScreenInteractor.swift
//  Fluently
//
//  Created by Савва Пономарев on 15.07.2025.
//

import Foundation
import SwiftUI

final class ProfileScreenInteractor: ObservableObject {
    let api: APIService

    init() {
        self.api = APIService()
    }

    func getPreferences() async throws -> PreferencesModel {
        return try await api.getPreferences()
    }

    func getGoals() async throws -> [String] {
        let response: [[String: String]] = try await api.getGoals()
        var topics: [String] = []
        response.forEach { item in
            if !topics.contains(item["title"]!) {
                topics.append(item["title"]!)
            }
        }
        topics.append("Learn new words")

        return topics
    }
}
