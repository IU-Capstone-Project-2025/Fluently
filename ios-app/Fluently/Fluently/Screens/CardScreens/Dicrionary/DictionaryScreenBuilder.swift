//
//  DictionaryScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

enum DictionaryScreenBuilder {
    static func build(
        isLearned: Bool
    ) -> DictionaryView {
        return DictionaryView(
            isLearned: isLearned
        )
    }
}
