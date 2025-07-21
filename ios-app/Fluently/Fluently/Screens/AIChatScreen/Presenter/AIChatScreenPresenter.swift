//
//  AIChatScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation
import SwiftUI

final class AIChatScreenPresenter: ObservableObject{

    var interactor: AIChatInteractor
#if targetEnvironment(simulator)
    @Published var messages: [MessageModel] = MessageModel.mockGenerator()
#else
    @Published var messages: [MessageModel] = []
#endif
    @Published var isReady = false

    init(interactor: AIChatInteractor) {
        self.interactor = interactor

        let initMessage = MessageModel(text: "", role: .user)
        Task {
            messages = try await interactor.sendMessage(chat: [initMessage])
        }
    }

    @MainActor
    func sendMessage(_ message: String) {
        let newMessage = MessageModel(text: message, role: .user)
        messages.append(newMessage)

        Task {
            messages = try await interactor.sendMessage(chat: messages)
        }

        if messages.count >= 10 {
            isReady = true
        }
    }

    func finishChat() {
        interactor.finishChat(chat: messages)
    }
}
