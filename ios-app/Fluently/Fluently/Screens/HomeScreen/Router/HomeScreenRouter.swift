//
//  HomeScreenRouter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

final class HomeScreenRouter {
    @ObservedObject var router: AppRouter

    init(router: AppRouter) {
        self.router = router
    }

    func navigatoToProfile() {
        router.navigate(to: AppRoutes.profile)
    }

    func navigatoToLesson(_ lesson: CardsModel) {
        router.navigate(to: AppRoutes.lesson(lesson))
    }
}
