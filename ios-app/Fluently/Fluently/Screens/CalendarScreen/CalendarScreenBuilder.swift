//
//  CalendarScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import Foundation

enum CalendarScreenBuilder {

    static func build(

    ) -> CalendarScreenView {
        let presenter = CalendarScreenPresenter()

        return CalendarScreenView(
            presenter: presenter
        )
    }
}
