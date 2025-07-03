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
            words: WordModel.generateMockWords(count: 5)
        )
        
        return LessonScreensView (
            presenter: presenter
        )
    }
}
