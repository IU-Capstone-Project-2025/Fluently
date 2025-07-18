//
//  MessageModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation
import SwiftData

//@Model
final class MessageModel: Codable {
    var text: String
    var role: MessageRole

    init(text: String, role: MessageRole) {
        self.text = text
        self.role = role
    }

    // MARK: - Codable

    enum CodingKeys: String, CodingKey {
        case text = "message"
        case role = "author"
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        text = try container.decode(String.self, forKey: .text)
        let roleString = try container.decode(String.self, forKey: .role)
        role = MessageRole(rawValue: roleString) ?? .ai
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(text, forKey: .text)
        try container.encode(role.rawValue, forKey: .role)
    }
}

extension MessageModel {
    static func mockGenerator() -> [MessageModel] {
        let mockMessages = [
            MessageModel(text: "Hello! How can I assist you today?", role: .ai),
            MessageModel(text: "Hi! I'm having trouble with my account.", role: .user),
            MessageModel(text: "I'm sorry to hear that. Could you please tell me what issue you're experiencing?", role: .ai),
            MessageModel(text: "I can't log in. It says my password is incorrect.", role: .user),
            MessageModel(text: "Have you tried resetting your password using the 'Forgot Password' option?", role: .ai),
            MessageModel(text: "Yes, but I didn't receive the reset email.", role: .user),
            MessageModel(text: "Let me check that for you. Could you confirm the email address associated with your account?", role: .ai),
            MessageModel(text: "It's example@email.com", role: .user),
            MessageModel(text: "Thank you. I see the issue - there was a typo in your email address. Would you like me to correct it?", role: .ai),
            MessageModel(text: "Yes, please! That would explain why I wasn't getting emails.", role: .user),
            MessageModel(text: "I've updated your email address. You should receive the password reset email shortly.", role: .ai),
            MessageModel(text: "Great, I got it! Let me try resetting now.", role: .user),
            MessageModel(text: "Perfect! Let me know if you need any further assistance.", role: .ai),
            MessageModel(text: "It worked! Thanks for your help.", role: .user),
            MessageModel(text: "You're welcome! Have a great day.", role: .ai)
        ]

        return mockMessages
    }
}
