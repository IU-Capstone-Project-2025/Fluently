//
//  APIService.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

protocol APIServiceProtocol {
    func authGoogle(_ gid: String) async throws -> AuthResponse
}

final class APIService: APIServiceProtocol{
    private let link = "https://fluently-app.ru"

    func authGoogle(_ gid: String) async throws -> AuthResponse {
        // Validate url
        guard let url = URL(string: link) else {
            throw ApiError.invalidURL
        }
        // Updating url
        let authGoogleUrl = url.appendingPathComponent("/auth/google")

        // Formating request
        var request = URLRequest(url: authGoogleUrl)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpMethod = "POST"

        // Add Data
        let requestBody = AuthRequest(idToken: gid)
        do {
            request.httpBody = try JSONEncoder().encode(requestBody)
        } catch {
            throw ApiError.encodingFailed
        }

        // Sending request
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse,
                  (200...299).contains(httpResponse.statusCode) else {
                throw ApiError.invalidResponse
            }

            let decoder = JSONDecoder()
            return try decoder.decode(AuthResponse.self, from: data)
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }

    func updateAccessToken() async throws {
        // Validate url
        guard let url = URL(string: link) else {
            throw ApiError.invalidURL
        }

        // Validate Refresh Token
        guard let refreshToken = KeyChainManager.shared.getRefreshToken() else {
            throw KeyChainManager.KeychainError.emptyRefreshToken
        }

        // Form request
        let refreshURL = url.appendingPathComponent("/auth/refresh")

        var request = URLRequest(url: refreshURL)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpMethod = "POST"

        let requestHttpBody = ["refresh_token" : refreshToken]
        do {
            request.httpBody = try JSONEncoder().encode(requestHttpBody)
        } catch {
            throw ApiError.encodingFailed
        }

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse,
                  (200...299).contains(httpResponse.statusCode) else {
                throw ApiError.invalidResponse
            }

            let decoder = JSONDecoder()
            let tokens = try decoder.decode(AuthResponse.self, from: data)

            // Try to upd keychain data
            do {
                try KeyChainManager.shared.saveToken(tokens)
            } catch {
                throw KeyChainManager.KeychainError.saveTokens
            }
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }
}

// MARK: - Lessons
extension APIService {
    func getLessons() async throws {
        // Validate url
        guard let url = URL(string: link) else {
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
        let lessonURL = url.appendingPathComponent("/lesson")

        var request = URLRequest(url: lessonURL)
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpMethod = "GET"

        let requestHttpBody = ["access_token" : accessToken]
        do {
            request.httpBody = try JSONEncoder().encode(requestHttpBody)
        } catch {
            throw ApiError.encodingFailed
        }

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse,
                  (200...299).contains(httpResponse.statusCode) else {
                throw ApiError.invalidResponse
            }

            let decoder = JSONDecoder()
            let decodedData = try decoder.decode(, from: data)
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }
}

extension APIService {
    // MARK: - Error
    enum ApiError: Error, Equatable {
        case invalidURL
        case encodingFailed
        case invalidResponse
        case networkError(String)
    }
}
