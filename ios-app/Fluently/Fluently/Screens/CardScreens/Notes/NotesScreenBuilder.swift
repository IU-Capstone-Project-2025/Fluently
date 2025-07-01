//
//  NotesScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

enum NotesScreenBuilder {

    static func build (
        router: AppRouter
    ) -> NotesView {
        let router = NotesScreenRouter(router: router)
        let presenter = NotesScreenPresenter(router: router)

        return NotesView(presenter: presenter)
    }
}
