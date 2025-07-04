//
//  HomeScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

// MARK: - Protocol for presenter
protocol HomeScreenPresenting: ObservableObject {
    func getLesson()

    // Navigation
    func navigatoToProfile()
    func navigatoToLesson()
}

// MARK: - Presenter implementation
final class HomeScreenPresenter: HomeScreenPresenting {
    let router: HomeScreenRouter
    let interactor: HomeScreenInteractor

    @ObservedObject var account: AccountData

    init(
        router: HomeScreenRouter,
        interactor: HomeScreenInteractor,
        account: AccountData
    ) {
        self.router = router
        self.interactor = interactor
        self.account = account
    }

    func getLesson() {
        interactor.getLesson()
    }

    // Builders 
    func buildNotesScreen() -> NotesView{
        return NotesScreenBuilder.build(router: router.router)
    }

    func buildDictionaryScreen() -> DictionaryView{
        return DictionaryScreenBuilder.build()
    }

    // Navigation

    func navigatoToProfile() {
        router.navigatoToProfile()
    }

    func navigatoToLesson() {
        router.navigatoToLesson()
    }
}
