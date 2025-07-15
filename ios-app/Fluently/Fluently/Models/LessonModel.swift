//
//  LessonModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import SwiftData

@Model
final class LessonModel: Codable {
    var startedAt: String
    var wordsPerLesson: Int
    var totalWords: Int
    var cefrLevel: String

    init(
        startedAt: String,
        totalWords: Int,
        wordsPerLesson: Int,
        cefrLevel: String
    ) {
        self.startedAt = startedAt
        self.totalWords = totalWords
        self.wordsPerLesson = wordsPerLesson
        self.cefrLevel = cefrLevel
    }

    // MARK: - Codable
    enum CodingKeys: String, CodingKey {
        case startedAt = "started_at"
        case totalWords = "total_words"
        case wordsPerLesson = "words_per_lesson"
        case cefrLevel = "cefr_level"
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        startedAt = try container.decode(String.self, forKey: .startedAt)
        totalWords = try container.decode(Int.self, forKey: .totalWords)
        wordsPerLesson = try container.decode(Int.self, forKey: .wordsPerLesson)
        cefrLevel = try container.decode(String.self, forKey: .cefrLevel)
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(startedAt, forKey: .startedAt)
        try container.encode(totalWords, forKey: .totalWords)
        try container.encode(wordsPerLesson, forKey: .wordsPerLesson)
        try container.encode(cefrLevel, forKey: .cefrLevel)
    }
}
