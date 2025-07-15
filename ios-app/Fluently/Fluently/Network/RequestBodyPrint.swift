//
//  RequestBodyPrint.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

// MARK: - URLRequests Extension
extension URLRequest {
    /// print the current HTTP Request 
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
