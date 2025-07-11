//
//  SentenceModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation
import SwiftData

@Model
final class SentenceModel: Codable{
    var id = UUID() 
    var text: String
    var translation: String

    init(
        text: String,
        translation: String
    ) {
        self.text = text
        self.translation = translation
    }

    enum CodingKeys: String, CodingKey {
        case text
        case translation
    }

    required init(from decoder: any Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)

        text = try container.decode(String.self, forKey: .text)
        translation = try container.decode(String.self, forKey: .translation)
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.container(keyedBy: CodingKeys.self)

        try container.encode(text, forKey: .text)
        try container.encode(translation, forKey: .translation)
    }
}

extension SentenceModel: Hashable {

    func hash(into hasher: inout Hasher) {
        hasher.combine(text)
        hasher.combine(translation)
    }

    static func == (lhs: SentenceModel, rhs: SentenceModel) -> Bool {
        return lhs.text == rhs.text && lhs.translation == rhs.translation
    }
}
