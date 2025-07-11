//
//  ApiServiceUnitTest.swift
//  FluentlyTests
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Testing
import GoogleSignIn
@testable import Fluently

struct ApiServiceUnitTest {

    @Test("Invalid URL Test")
    func namehandleInvalidURLError() async throws {
        // Given
        let mockAPIService = MockAPIService()
        mockAPIService.shouldSucceed = false

        // Then
        await #expect(throws: APIService.ApiError.invalidURL) {
            _ = try await mockAPIService.authGoogle("")
        }
    }

    @Test("Invalid Response Test")
    func namehandleInvalidResponseError() async throws {
        // Given
        let mockAPIService = MockAPIService()
        mockAPIService.shouldSucceed = false

        // Then
        await #expect(throws: APIService.ApiError.invalidResponse(statusCode: -1)) {
            _ = try await mockAPIService.authGoogle("123")
        }
    }

    @Test("Normal auth")
    func authDataHandle() async throws {
        // Given
        let mockAPIService = MockAPIService()
        mockAPIService.shouldSucceed = true

        // When
        let response = try await mockAPIService.authGoogle("123")

        // Then
        #expect(response.accessToken.isEmpty == false)
        #expect(response.refreshToken.isEmpty == false)
        #expect(response.tokenType.isEmpty == false)
    }
}

class MockAPIService: AuthAPI{
    var shouldSucceed = true

    func updateAccessToken() async throws {
        if shouldSucceed {
            _ = AuthResponse(
                accessToken: "mock_access_token",
                refreshToken: "mock_refresh_token",
                tokenType: "Bearer",
                expiresIn: 3600
            )
        } else {
            throw APIService.ApiError.invalidResponse(statusCode: -1)
        }
    }

    func authGoogle(_ gid: String) async throws -> AuthResponse {
        guard gid.isEmpty == false else {
            throw APIService.ApiError.invalidURL
        }

        if shouldSucceed {
            return AuthResponse(
                accessToken: "mock_access_token",
                refreshToken: "mock_refresh_token",
                tokenType: "Bearer",
                expiresIn: 3600
            )
        } else {
            throw APIService.ApiError.invalidResponse(statusCode: -1)
        }
    }
}
