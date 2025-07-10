//
//  LessonModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

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

    enum CodingKeys: String, CodingKey {
        case startedAt = "started_at"
        case totalWords = "total_words"
        case wordsPerLesson = "words_per_lesson"
        case cefrLevel = "cefr_level"
    }
}
