//
//  LessonScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation

enum LessonScreenBuilder {
    static func build (
        router: AppRouter,
        lesson: [WordModel]
    ) -> LessonScreensView {
        let presenter = LessonsPresenter (
            router: router,
            words: lesson
        )

        return LessonScreensView (
            presenter: presenter
        )
    }
}
