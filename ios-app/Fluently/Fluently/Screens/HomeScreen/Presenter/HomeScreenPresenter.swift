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

    // Nacigation
    func navigatoToProfile()
}

// MARK: - Presenter implementation
final class HomeScreenPresenter: HomeScreenPresenting {
    @ObservedObject var router: AppRouter
    @ObservedObject var account: AccountData

    init(router: AppRouter, account: AccountData) {
        self.router = router
        self.account = account
    }

    func navigatoToProfile() {
        router.navigate(to: AppRoutes.profile)
    }
}
