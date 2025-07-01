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

    // Navigation
    func navigatoToProfile()
    func navigatoToLesson()
}

// MARK: - Presenter implementation
final class HomeScreenPresenter: HomeScreenPresenting {
    let router: HomeScreenRouter

    @ObservedObject var account: AccountData

    init(
        router: HomeScreenRouter,
        account: AccountData
    ) {
        self.router = router
        self.account = account
    }

    func navigatoToProfile() {
        router.navigatoToProfile()
    }

    func navigatoToLesson() {
        router.navigatoToLesson()
    }
}
