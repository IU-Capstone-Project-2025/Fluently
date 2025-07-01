//
//  ApiError.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

// MARK: - Error
enum ApiError: Error {
    case invalidURL
    case encodingFailed
    case invalidResponse
    case networkError(String)
}
