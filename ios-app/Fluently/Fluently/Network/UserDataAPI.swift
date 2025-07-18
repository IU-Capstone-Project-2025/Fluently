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
}
