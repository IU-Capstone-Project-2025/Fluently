//
//  APIService.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

final class APIService{
    private let baseUrl = "https://fluently-app.ru"

    func getURL() -> String {
        return baseUrl
    }

    func sendRequest(request: URLRequest) async throws -> Data  {
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse,
                  (200...299).contains(httpResponse.statusCode) else {
                throw ApiError.invalidResponse
            }

            return data
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }

    func getTokens() async throws -> AuthResponse {
        // Validate url
        guard let url = URL(string: baseUrl) else {
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
            let data = try await sendRequest(request: request)

            let decoder = JSONDecoder()
            let tokens = try decoder.decode(AuthResponse.self, from: data)
            return tokens
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }
}


// MARK: - Error
extension APIService {
    enum ApiError: Error, Equatable {
        case invalidURL
        case encodingFailed
        case invalidResponse
        case networkError(String)

        var localizedDescription: String {
            switch self {
               case .invalidURL: return "Invalid URL"
               case .encodingFailed: return "Failed to encode request data"
               case .invalidResponse: return "Received invalid response from server"
               case .networkError(let error): return "Network error: \(error)"
            }
        }
    }
}
