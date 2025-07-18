//
//  ProfileScreenInteractor.swift
//  Fluently
//
//  Created by Савва Пономарев on 15.07.2025.
//

import Foundation
import SwiftUI

final class ProfileScreenInteractor: ObservableObject {
    let api: APIService

    init() {
        self.api = APIService()
    }

    func getPreferences() async throws -> PreferencesModel {
        return try await api.getPreferences()
    }
}
