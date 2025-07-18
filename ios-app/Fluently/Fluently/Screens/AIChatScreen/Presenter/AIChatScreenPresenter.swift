//
//  AIChatScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation
import SwiftUI

final class AIChatScreenPresenter: ObservableObject{
    @Published var messages: [MessageModel] = MessageModel.mockGenerator()
    @Published var isReady = false

    func sendMessage(_ message: String) {
        let newMessage = MessageModel(text: message, role: .user)
        messages.append(newMessage)

        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            let aiResponse = MessageModel(text: "This is an automated response", role: .ai)
            self.messages.append(aiResponse)
        }

        if messages.count >= 10 {
            isReady = true
        }
    }
}
