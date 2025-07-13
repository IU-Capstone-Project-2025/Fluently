//
//  APIService.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation


final class APIService {
    let baseUrl = "https://fluently-app.ru"

    func sendRequest(_ request: URLRequest) async throws -> Data {
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            request.printRequest()

            if let jsonString = String(data: data, encoding: .utf8) {
                print("Request: \(request.url?.absoluteString ?? "")")
                print("Raw JSON Response:\n\(jsonString)")
            }

            guard let httpResponse = response as? HTTPURLResponse else {
                throw ApiError.invalidResponse(statusCode: nil)
            }

            print(httpResponse.statusCode)

            guard (200...299).contains(httpResponse.statusCode) else {
                throw ApiError.invalidResponse(statusCode: httpResponse.statusCode)
            }

            return data
        } catch let error as URLError {
            throw ApiError.networkError(
                code: error.errorCode,
                message: error.localizedDescription
            )
        } catch {
            throw ApiError.networkError(
                code: -1,
                message: error.localizedDescription
            )
        }
    }

    func makeRequest<T: Encodable>(
        path: String,
        method: String,
        body: T? = nil,
        headers: [String: String] = [:]
    ) throws -> URLRequest {
        guard let url = URL(string: baseUrl)?.appendingPathComponent(path) else {
            throw ApiError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        headers.forEach { request.setValue($1, forHTTPHeaderField: $0) }

        if let body = body {
            request.httpBody = try JSONEncoder().encode(body)
        }

        return request
    }

    func decodeResponse<T: Decodable>(from request: URLRequest) async throws -> T {
        let data = try await sendRequest(request)
        do {
            return try JSONDecoder().decode(T.self, from: data)
        } catch let error as DecodingError {
            print("JSON Decoding Error: \(error.localizedDescription)")
            switch error {
                case .typeMismatch(let type, let context):
                    print("Type mismatch for \(type): \(context.debugDescription)")
                case .valueNotFound(let type, let context):
                    print("Value not found for \(type): \(context.debugDescription)")
                case .keyNotFound(let key, let context):
                    print("Key '\(key.stringValue)' not found: \(context.debugDescription)")
                case .dataCorrupted(let context):
                    print("Data corrupted: \(context.debugDescription)")
                @unknown default:
                    print("Unknown error: \(error)")
            }
            throw ApiError.decodingFailed(error.localizedDescription)
        }
    }
}


extension APIService {
    // MARK: - Error
    enum ApiError: Error, Equatable {
        case invalidURL
        case encodingFailed (String)
        case decodingFailed (String)
        case invalidResponse (statusCode: Int?)
        case networkError (code: Int, message: String)
    }
}
