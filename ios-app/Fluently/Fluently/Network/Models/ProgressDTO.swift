//
//  ProgressDTO.swift
//  Fluently
//
//  Created by Савва Пономарев on 13.07.2025.
//

import Foundation

// MARK: - Encoding form of progress
/// used to send progress back to the backend
struct ProgressDTO : Encodable {
    var cnt_reviewed: Int = 1
    var confidence_score: Int = 100
    var learned_at: String
    var word_id: String

    init(word_id: String) {
        self.cnt_reviewed = 1
        self.confidence_score = 100
        self.learned_at = Date.now.ISO8601Format()
        self.word_id = word_id
    }

    enum CodingKeys: String, CodingKey {
        case cnt_reviewed
        case confidence_score
        case learned_at
        case word_id
    }

    func encode(to encoder: any Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(cnt_reviewed, forKey: .cnt_reviewed)
        try container.encode(confidence_score, forKey: .confidence_score)
        try container.encode(learned_at, forKey: .learned_at)
        try container.encode(word_id, forKey: .word_id)
    }
}
