//
//  LessonAPIProtocol.swift
//  Fluently
//
//  Created by Савва Пономарев on 04.07.2025.
//

import Foundation

protocol LessonAPIProtocol {
    func getLessons() async throws
}

// MARK: - Lessons
extension APIService: LessonAPIProtocol{
    func getLessons() async throws {
        // Validate url
        guard let url = URL(string: getURL()) else {
            throw ApiError.invalidURL
        }

        // Validate Refresh Token
        if !KeyChainManager.shared.isTokenValid() {
            try await updateAccessToken()
        }

        guard let accessToken = KeyChainManager.shared.getAccessToken() else {
            throw KeyChainManager.KeychainError.emptyAccessToken
        }

        // Form request
        let lessonURL = url.appendingPathComponent("/api/v1/lesson")

        var request = URLRequest(url: lessonURL)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        request.httpMethod = "GET"

        request.printRequest()

        do {
            let data = try await sendRequest(request: request)

            let decoder = JSONDecoder()
            decoder.keyDecodingStrategy = .convertFromSnakeCase
            let decodedData = try decoder.decode(CardsModel.self, from: data)

            print(decodedData)
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }
}
