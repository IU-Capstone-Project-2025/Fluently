//
//  KeyChainManager.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import Security

final class KeyChainManager {
    private init() { }

    static let shared = KeyChainManager()

    private let service = "BRPM.Fluently"

    // MARK: - Save Token
    func saveToken(_ response: AuthResponse) throws {
        try save(response.accessToken, service: service, account: "accessToken")
        try save(response.refreshToken, service: service, account: "refreshToken")
        UserDefaults.standard.set(Date().addingTimeInterval(TimeInterval(response.expiresIn)), forKey: "tokenExpiry")
    }

    // MARK: - Retrieve Tokens
    func getAccessToken() -> String? {
        try? get(service: service, account: "accessToken")
    }

    func getRefreshToken() -> String? {
        try? get(service: service, account: "refreshToken")
    }

    func isTokenValid() -> Bool {
        guard let expiryDate = UserDefaults.standard.object(forKey: "tokenExpiry") as? Date else {
            return false
        }
        return expiryDate > Date()
    }

    // MARK: - Delete Tokens
    func deleteTokens() throws {
        try delete(service: service, account: "accessToken")
        try delete(service: service, account: "refreshToken")
        UserDefaults.standard.removeObject(forKey: "tokenExpiry")
    }

}

// MARK: - Private
private extension KeyChainManager {

    func save(_ value: String, service: String, account: String) throws {
        guard let data = value.data(using: .utf8) else {
            throw KeychainError.encodingError
        }

        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: account,
            kSecValueData: data,
            kSecAttrAccessible: kSecAttrAccessibleAfterFirstUnlock
        ]

        SecItemDelete(query as CFDictionary)
        let status = SecItemAdd(query as CFDictionary, nil)

        guard status == errSecSuccess else {
            throw KeychainError.unhandledError(status: status)
        }
    }

    func get(service: String, account: String) throws -> String {
        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: account,
            kSecReturnData: kCFBooleanTrue!,
            kSecMatchLimit: kSecMatchLimitOne
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess,
              let data = result as? Data,
              let value = String(data: data, encoding: .utf8) else {
            throw KeychainError.noTokenFound
        }

        return value
    }

    func delete(service: String, account: String) throws {
        let query: [CFString: Any] = [
            kSecClass: kSecClassGenericPassword,
            kSecAttrService: service,
            kSecAttrAccount: account
        ]

        let status = SecItemDelete(query as CFDictionary)
        guard status == errSecSuccess || status == errSecItemNotFound else {
            throw KeychainError.unhandledError(status: status)
        }
    }
}

// MARK: - Keychain Errors
extension KeyChainManager {
    enum KeychainError: Error {
        case encodingError
        case noTokenFound
        case unhandledError(status: OSStatus)

        case emptyRefreshToken
        case saveTokens

        var localizedDescription: String {
            switch self {
                case .encodingError: return "Failed to encode token"
                case .noTokenFound: return "Token not found in Keychain"
                case .unhandledError(let status):
                    return "Keychain error: \(status)"
                case .emptyRefreshToken: return "The Refresh token in empty"
                case .saveTokens: return "Failed save tokens"
            }
        }
    }
}
