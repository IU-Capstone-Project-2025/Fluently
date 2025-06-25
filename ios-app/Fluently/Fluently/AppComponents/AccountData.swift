//
//  AccountData.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.06.2025.
//

import SwiftUI

final class AccountData: ObservableObject {
    @Published var name: String?
    @Published var familyName: String?
    @Published var mail: String?
    @Published var image: String?
    @Published var isLoggined = false

    let defaults = UserDefaults.standard

    func saveData() async {
        let encoder = JSONEncoder()
        if let logstatus = try? encoder.encode(isLoggined) {
            defaults.set(logstatus, forKey: "isLoggined")
        }
    }

    func read() throws {
        if let savedData = defaults.object(forKey: "isLoggined") as? Data {
            let decoder = JSONDecoder()
            do {
                let isLoggined = try? decoder.decode(Bool.self, from: savedData)
            } catch {
                throw AccountError.decodingError
            }
        }
    }
}

enum AccountError: Error {
    case decodingError
}
