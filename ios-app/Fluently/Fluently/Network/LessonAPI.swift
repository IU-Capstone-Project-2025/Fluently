//
//  LessonAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 05.07.2025.
//

import Foundation

protocol LessonAPI {
    func getLesson() async throws
}

// MARK: - Lessons
extension APIService: LessonAPI {
    func getLesson() async throws {
        // Validate Refresh Token
        try await validateToken()

        let request = try makeAuthorizedRequest(
            path: "/api/v1/lesson",
            method: "GET"
        )

        print(try await fetchAndDecode(request: request) as CardsModel)
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
            "Bearer \(accessToken)",
            forHTTPHeaderField: "Authorization"
        )

        return request
    }

    private func fetchAndDecode<T: Decodable>(
        request: URLRequest,
        decoder: JSONDecoder = JSONDecoder()
    ) async throws -> T {
        let data = try await sendRequest(request)
        decoder.keyDecodingStrategy = .convertFromSnakeCase

        do {
            return try decoder.decode(T.self, from: data)
        } catch let error as DecodingError {
            throw ApiError.decodingFailed(error.localizedDescription)
        }
    }
}
