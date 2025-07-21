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
            method: "GET",
            body: Optional<String>.none
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

//        let progressItems = words.map { word in
//            ProgressDTO(
//                word_id: word.wordId ?? UUID().uuidString
//            )
//        }

        let progressItems: [WordProgressItem] = words.compactMap { word in
            guard let wordId = word.wordId else { return nil }
            return word.isLearned ?
                .learned(ProgressDTO(word_id: wordId)) :
                .notLearned(wordId)
        }


        var request = try makeRequest(
            path: path,
            method: method,
            body: progressItems
        )

        request.setValue(
            "Bearer \(accessToken)", forHTTPHeaderField: "Authorization"
        )

        let data = try await sendRequest(request)
        print(String(data: data, encoding: .utf8) ?? "nil data returned")
    }

    func getDayWord() async throws -> WordModel {
        try await validateToken()

        let path = "/api/v1/day-word"
        let method = "GET"

        let request = try makeAuthorizedRequest(
            path: path,
            method: method,
            body: Optional<String>.none
        )

        return try await fetchAndDecode(request: request)
    }
}

enum WordProgressItem: Encodable {
    case learned(ProgressDTO)
    case notLearned(String)

    enum CodingKeys: String, CodingKey {
        case wordId = "word_id"
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)
        switch self {
            case .learned(let dto):
                try dto.encode(to: encoder)
            case .notLearned(let id):
                try container.encode(id, forKey: .wordId)
        }
    }
}
