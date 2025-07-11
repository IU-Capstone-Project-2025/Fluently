//
//  StatisticScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import Foundation


enum StatisticScreenBuilder {

    static func build(

    ) -> StatisticScreenView {
        let presenter = StatisticScreenPresenter()

        return StatisticScreenView(presenter: presenter)
    }
}
