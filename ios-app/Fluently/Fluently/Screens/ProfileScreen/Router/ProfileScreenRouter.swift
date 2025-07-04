//
//  ProfileScreenRouter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import SwiftUI

final class ProfileScreenRouter {
    @ObservedObject var router: AppRouter

    init(router: AppRouter) {
        self.router = router
    }

    func popToRoot() {
        router.popToRoot()
    }

    func navigateBack() {
        router.pop()
    }

    func navigateToLogin() {
        router.navigate(to: AppRoutes.login)
    }
}
