//
//  LessonScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

enum LessonScreenBuilder {

    static func build (
        router: AppRouter
    ) -> LessonScreensView {
        let presenter = LessonsPresenter (
            router: router,
            words: WordCardGenerator.generateCards()
        )
        
        return LessonScreensView (
            presenter: presenter
        )
    }
}
