//
//  Token.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

// MARK: - Request Struct
struct AuthRequest: Codable {
    let idToken: String
    let platform: String = "iOS"

    enum CodingKeys: String, CodingKey {
        case idToken = "id_token"
        case platform
    }
}

// MARK: - Response Struct
struct AuthResponse: Codable {
    let accessToken: String
    let refreshToken: String
    let tokenType: String
    let expiresIn: Int

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
        case tokenType = "token_type"
        case expiresIn = "expires_in"
    }
}
