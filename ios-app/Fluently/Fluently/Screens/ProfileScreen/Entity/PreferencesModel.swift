//
//  PreferencesModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 15.07.2025.
//

import Foundation
import SwiftData

@Model
final class PreferencesModel: Codable, Sendable{
    var avatarURL: String
    var cefrLevel: String
    var dailyWord: Bool
    var goal: String
    var id: String      /// `id` for preferences
    var notificationAt: Date
    var notifications: Bool
    var subscribed: Bool
    var userId: String  /// `id` of the user
    var wordPerDay: Int

    init(
        avatarURL: String,
        cefrLevel: String,
        dailyWord: Bool,
        goal: String,
        id: String,
        notificationAt: Date,
        notifications: Bool,
        subscribed: Bool,
        userId: String,
        wordPerDay: Int
    ) {
        self.avatarURL = avatarURL
        self.cefrLevel = cefrLevel
        self.dailyWord = dailyWord
        self.goal = goal
        self.id = id
        self.notificationAt = notificationAt
        self.notifications = notifications
        self.subscribed = subscribed
        self.userId = userId
        self.wordPerDay = wordPerDay
    }

    // MARK: - Decodable

    enum CodingKeys: String, CodingKey {
        case avatarURL = "avatar_image_url"
        case cefrLevel = "cefr_level"
        case dailyWord = "fact_everyday"
        case goal
        case id
        case notificationAt = "notification_at"
        case notifications
        case subscribed
        case userId = "user_id"
        case wordPerDay = "words_per_day"
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        self.avatarURL = try container.decode(String.self, forKey: .avatarURL)
        self.cefrLevel = try container.decode(String.self, forKey: .cefrLevel)
        self.dailyWord = try container.decode(Bool.self, forKey: .dailyWord)
        self.goal = try container.decode(String.self, forKey: .goal)
        self.id = try container.decode(String.self, forKey: .id)
        let notificationDateISO = try container.decodeIfPresent(String.self, forKey: .notificationAt) ?? Date.now.ISO8601Format()

        let dateFormatter = DateFormatter()
        dateFormatter.locale = Locale(identifier: "en_US_POSIX")
        dateFormatter.dateFormat = "yyyy-MM-dd'T'HH:mm:ssZ"
        self.notificationAt = dateFormatter.date(from:notificationDateISO)!

        self.notifications = try container.decode(Bool.self, forKey: .notifications)
        self.subscribed = try container.decode(Bool.self, forKey: .subscribed)
        self.userId = try container.decode(String.self, forKey: .userId)
        self.wordPerDay =  try container.decode(Int.self, forKey: .wordPerDay)
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(avatarURL, forKey: .avatarURL)
        try container.encode(cefrLevel, forKey: .cefrLevel)
        try container.encode(dailyWord, forKey: .dailyWord)
        try container.encode(goal, forKey: .goal)
        try container.encode(id, forKey: .id)
        try container.encode(notificationAt.ISO8601Format(), forKey: .notificationAt)
        try container.encode(notifications, forKey: .notifications)
        try container.encode(subscribed, forKey: .subscribed)
        try container.encode(userId, forKey: .userId)
        try container.encode(wordPerDay, forKey: .wordPerDay)
    }
}


extension PreferencesModel {
    static func generate(
        avatarURL: String? = nil,
        cefrLevel: String? = nil,
        dailyWord: Bool? = nil,
        goal: String? = nil,
        id: String? = nil,
        notificationAt: String? = nil,
        notifications: Bool? = nil,
        subscribed: Bool? = nil,
        userId: String? = nil,
        wordPerDay: Int? = nil,
        randomize: Bool = false
    ) -> PreferencesModel {

        let randomAvatar = [
            "https://example.com/avatars/user1.jpg",
            "https://example.com/avatars/user2.png",
            "https://example.com/avatars/default_avatar.svg"
        ].randomElement()!

        let randomCefr = ["A1", "A2", "B1", "B2", "C1", "C2"].randomElement()!
        let randomGoals = ["fluency", "travel", "business", "exams", "culture"]

        return PreferencesModel(
            avatarURL: avatarURL ?? (randomize ? randomAvatar : "https://example.com/avatars/default.jpg"),
            cefrLevel: cefrLevel ?? (randomize ? randomCefr : "B2"),
            dailyWord: dailyWord ?? (randomize ? Bool.random() : true),
            goal: goal ?? (randomize ? randomGoals.randomElement()! : "fluency"),
            id: id ?? (randomize ? UUID().uuidString : "pref_12345"),
            notificationAt: Date.now,
            notifications: notifications ?? (randomize ? Bool.random() : true),
            subscribed: subscribed ?? (randomize ? Bool.random() : false),
            userId: userId ?? (randomize ? UUID().uuidString : "user_67890"),
            wordPerDay: wordPerDay ?? (randomize ? Int.random(in: 5...50) : 10)
            )
    }
}
