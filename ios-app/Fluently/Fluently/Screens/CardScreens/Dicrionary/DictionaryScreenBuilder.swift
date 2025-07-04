//
//  DictionaryScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

enum DictionaryScreenBuilder {
    static func build(

    ) -> DictionaryView {
        let presenter = DictionaryScreenPresenter()

        return DictionaryView(
            presenter: presenter
        )
    }
}
