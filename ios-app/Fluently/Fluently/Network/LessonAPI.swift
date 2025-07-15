//
//  LessonAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 05.07.2025.
//

import Foundation

// MARK: - Protocol
protocol LessonAPI {
    func getLesson() async throws -> CardsModel
    func sendProgress(words: [WordModel]) async throws
    func getDayWord() async throws -> WordModel
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

    func sendProgress(words: [WordModel]) async throws {
        try await validateToken()

        guard let accessToken = KeyChainManager.shared.getAccessToken() else {
            throw KeyChainManager.KeychainError.emptyAccessToken
        }

        let path = "/api/v1/progress"
        let method = "POST"

        let progressItems = words.map { word in
            ProgressDTO(
                word_id: word.wordId
            )
        }

        var request = try makeRequest(
            path: path,
            method: method,
            body: progressItems
        )

        request.setValue(
            "Bearer \(accessToken)", forHTTPHeaderField: "Authorization"
        )

        _ = try await sendRequest(request)
    }

    func getDayWord() async throws -> WordModel {
        try await validateToken()

        let path = "/api/v1/day-word"
        let method = "GET"

        let request = try makeAuthorizedRequest(
            path: path,
            method: method
        )

        return try await fetchAndDecode(request: request)
    }
}
