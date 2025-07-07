//
//  SentenceModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

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
