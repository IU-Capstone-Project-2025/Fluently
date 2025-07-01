//
//  APIService.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation


final class APIService{
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
}

extension URLRequest {
    func printRequest() {
        print("=== HTTP Request ===")
        print("URL: \(self.httpMethod ?? "GET") \(self.url?.absoluteString ?? "No URL")")

        // Print Headers
        print("Headers:")
        self.allHTTPHeaderFields?.forEach { print("  \($0.key): \($0.value)") }

        // Print Body
        if let body = self.httpBody {
            print("Body:")
            if let json = try? JSONSerialization.jsonObject(with: body),
               let prettyData = try? JSONSerialization.data(withJSONObject: json, options: .prettyPrinted),
               let prettyString = String(data: prettyData, encoding: .utf8) {
                print(prettyString)
            } else {
                print(String(data: body, encoding: .utf8) ?? "Unable to decode body")
            }
        } else {
            print("Body: Empty")
        }
        print("===================")
    }
}
