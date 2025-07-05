//
//  AuthAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 04.07.2025.
//

import Foundation

protocol AuthAPIProtocol {
    func authGoogle(_ gid: String) async throws -> AuthResponse
    func updateAccessToken() async throws
}

// MARK: - Auth
extension APIService: AuthAPIProtocol {
    func authGoogle(_ gid: String) async throws -> AuthResponse {
        // Validate url
        guard let url = URL(string: getURL()) else {
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
            let data = try await sendRequest(request: request)

            let decoder = JSONDecoder()
            return try decoder.decode(AuthResponse.self, from: data)
        } catch {
            throw ApiError.networkError(error.localizedDescription)
        }
    }

    func updateAccessToken() async throws {
        let tokens = try await getTokens()

        // Try to upd keychain data
        do {
            try KeyChainManager.shared.saveToken(tokens)
        } catch {
            throw KeyChainManager.KeychainError.saveTokens
        }
    }
}
