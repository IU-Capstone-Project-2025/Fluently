//
//  AccountData.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.06.2025.
//

import SwiftUI

final class AccountData: ObservableObject {
    // MARK: - Properties
    @Published var name: String?
    @Published var familyName: String?
    @Published var mail: String?
    @Published var image: String?
    @Published var isLoggedIn = false
    
    // MARK: - Data caching
    let defaults = UserDefaults.standard

    func saveData() async {
        let encoder = JSONEncoder()
        if let logstatus = try? encoder.encode(isLoggedIn) {
            defaults.set(logstatus, forKey: "isLoggedIn")
        }
    }

    func read() throws {
        if let savedData = defaults.object(forKey: "isLoggedIn") as? Data {
            let decoder = JSONDecoder()
            do {
                _ = try? decoder.decode(Bool.self, from: savedData)
            }
        }
    }
}

// MARK: - Account Errors
enum AccountError: Error {
    case decodingError
}
