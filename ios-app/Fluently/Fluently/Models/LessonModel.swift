//
//  LessonModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

final class LessonModel: Codable {
    var cefrLevel: String
    var startedAt: String
    var totalWords: Int
    var wordsPerLesson: Int

    init(
        cefrLevel: String,
        startedAt: String,
        totalWords: Int,
        wordsPerLesson: Int
    ) {
        self.cefrLevel = cefrLevel
        self.startedAt = startedAt
        self.totalWords = totalWords
        self.wordsPerLesson = wordsPerLesson
    }

    enum CodingKeys: CodingKey {
        case cefrLevel
        case startedAt
        case totalWords
        case wordsPerLesson
    }
}
