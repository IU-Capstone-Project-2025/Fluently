//
//  AIChatAPI.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.07.2025.
//

import Foundation

protocol AIChatAPI {
    func sendMessage(chat: [MessageModel]) async throws -> [MessageModel]
    func finishChat(chat: [MessageModel]) async throws
    func getHistory() async throws -> [MessageModel]
}

extension APIService: AIChatAPI {
    func sendMessage(chat: [MessageModel]) async throws -> [MessageModel] {
        do {
            try await validateToken()

            let path = "api/v1/chat"
            let method = "POST"
            let body = ["chat": chat]

            let request = try makeAuthorizedRequest(
                path: path,
                method: method,
                body: body
            )

            let response: ChatResponse = try await fetchAndDecode(request: request)
            return response.chat
        } catch {
            print(error.localizedDescription)
            let message = MessageModel(text: "\(error.localizedDescription)", role: .ai)
            var newChat = chat
            newChat.append(message)
            return newChat
        }
    }
    
    func finishChat(chat: [MessageModel]) async throws {
        do {
            try await validateToken()

            let path = "api/v1/chat/finish"
            let method = "POST"
            let body = ["chat": chat]

            let request = try makeAuthorizedRequest(
                path: path,
                method: method,
                body: body
            )

            let _: EmptyResponse = try await fetchAndDecode(request: request)
        } catch {
            print(error.localizedDescription)
        }
    }
    
    func getHistory() async throws -> [MessageModel] {
        do {
            try await validateToken()

            let path = "api/v1/chat/finish"
            let method = "GET"

            let request = try makeAuthorizedRequest(
                path: path,
                method: method,
                body: Optional<String>.none
            )

            let response: [MessageModel] = try await fetchAndDecode(request: request)
            return response
        } catch {
            print(error.localizedDescription)
            return []
        }
    }

    struct ChatResponse: Decodable {
        let chat: [MessageModel]
    }

    struct EmptyResponse: Decodable {}
}
