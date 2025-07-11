//
//  LessonAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 05.07.2025.
//

import Foundation

protocol LessonAPI {
    func getLesson() async throws -> CardsModel
}

// MARK: - Lessons
extension APIService: LessonAPI {
    func getLesson() async throws -> CardsModel {
        // Validate Refresh Token
        try await validateToken()

        let request = try makeAuthorizedRequest(
            path: "/api/v1/lesson",
            method: "GET"
        )

        return try await fetchAndDecode(request: request) as CardsModel
    }

    // MARK: - Private Helpers

    private func validateToken() async throws {
        if !KeyChainManager.shared.isTokenValid() {
            try await updateAccessToken()
        }
    }

    private func makeAuthorizedRequest(
        path: String,
        method: String,
        headers: [String: String] = [:]
    ) throws -> URLRequest {
        guard let accessToken = KeyChainManager.shared.getAccessToken() else {
            throw KeyChainManager.KeychainError.emptyAccessToken
        }

        var request = try makeRequest(
            path: path,
            method: method,
            body: Optional<String>.none,
            headers: headers
        )

        request.setValue(
            "Bearer \(accessToken)", forHTTPHeaderField: "Authorization"
        )

        return request
    }

    private func fetchAndDecode<T: Decodable>(
        request: URLRequest,
        decoder: JSONDecoder = JSONDecoder()
    ) async throws -> T {
        let data = try await sendRequest(request)

        do {
            return try decoder.decode(T.self, from: data)

        } catch let error as DecodingError {
            print("JSON Decoding Error: \(error.localizedDescription)")
            switch error {
                case .typeMismatch(let type, let context):
                    print("Type mismatch for \(type): \(context.debugDescription)")
                case .valueNotFound(let type, let context):
                    print("Value not found for \(type): \(context.debugDescription)")
                case .keyNotFound(let key, let context):
                    print("Key '\(key.stringValue)' not found: \(context.debugDescription)")
                case .dataCorrupted(let context):
                    print("Data corrupted: \(context.debugDescription)")
                @unknown default:
                    print("Unknown error: \(error)")
            }
            throw ApiError.decodingFailed(error.localizedDescription)
        }
    }
}
