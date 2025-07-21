//
//  AuthAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 05.07.2025.
//

import Foundation

// MARK: - Protocol
protocol AuthAPI {
    // Tokens
    func authGoogle(_ gid: String) async throws -> AuthResponse
    func updateAccessToken() async throws
}

// MARK: - Auth Logic
extension APIService: AuthAPI {
    func authGoogle(_ gid: String) async throws -> AuthResponse {
        let request = try makeRequest(
            path: "/auth/google",
            method: "POST",
            body: AuthRequest(idToken: gid)
        )
        return try await decodeResponse(from: request)
    }

    func updateAccessToken() async throws {
        let tokens = try await refreshTokens()
        try KeyChainManager.shared.saveToken(tokens)
    }
}

// MARK: - Private
private extension APIService {
    func refreshTokens() async throws -> AuthResponse {
        guard let refreshToken = KeyChainManager.shared.getRefreshToken() else {
            throw KeyChainManager.KeychainError.emptyRefreshToken
        }

        let request = try makeRequest(
            path: "/auth/refresh",
            method: "POST",
            body: ["refresh_token": refreshToken]
        )
        return try await decodeResponse(from: request)
    }
}
