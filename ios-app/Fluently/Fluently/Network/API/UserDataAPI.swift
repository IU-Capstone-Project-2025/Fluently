//
//  UserDataAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 15.07.2025.
//

import Foundation

// MARK: - Protocol
protocol UserDataAPI {
    func getPreferences() async throws -> PreferencesModel
    func updatePreferences(_ prefs: PreferencesModel) async throws
    func getGoals() async throws -> [[String : String]]
}

// MARK: - User Data
extension APIService: UserDataAPI {
    func getPreferences() async throws -> PreferencesModel {
        try await validateToken()

        let path = "api/v1/preferences"
        let method = "GET"

        let request = try makeAuthorizedRequest(
            path: path,
            method: method,
            body: Optional<String>.none
        )

        return try await fetchAndDecode(request: request)
    }

    func updatePreferences(_ prefs: PreferencesModel) async throws {
        try await validateToken()

        let path = "api/v1/preferences"
        let method = "PUT"

        let request = try makeAuthorizedRequest(
            path: path,
            method: method,
            body: prefs
        )

        let _: Data = try await fetchAndDecode(request: request)
    }

    func getGoals() async throws -> [[String : String]] {
        try await validateToken()

        let path = "api/v1/topics"
        let method = "GET"

        let request = try makeAuthorizedRequest(
            path: path,
            method: method,
            body: Optional<String>.none
        )

        return try await fetchAndDecode(request: request)
    }
}
