//
//  AIChatInteractor.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.07.2025.
//

import Foundation

final class AIChatInteractor: ObservableObject {
    var api: APIService

    init() {
        self.api = APIService()
    }
    
    func sendMessage(chat: [MessageModel]) async throws -> [MessageModel] {
        return try await api.sendMessage(chat: chat)
    }

    func finishChat(chat: [MessageModel]) {
        Task {
            try await api.finishChat(chat: chat)
        }
    }
}
