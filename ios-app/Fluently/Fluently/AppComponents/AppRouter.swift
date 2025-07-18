//
//  AppRouter.swift
//  Fluently
//
//  Created by Савва Пономарев on 23.06.2025.
//

import Foundation
import SwiftUI

final class AppRouter: ObservableObject {
    @Published var navigationPath: NavigationPath

    init(initialPath: NavigationPath = NavigationPath()) {
        self.navigationPath = initialPath
    }

    /// returns the current Navigation path
    func getPath() -> NavigationPath {
        return navigationPath
    }

    /// navigate to sprecific screen basing on
    /// `destination` -> `AppRoutes`
    func navigate(to destination: any Hashable) {
        navigationPath.append(destination)
    }

    /// returns to the first screen
    func popToRoot() {
        navigationPath.removeLast(navigationPath.count)
    }

    /// returns to previous screen
    func pop() {
        if navigationPath.count >= 1  {
            navigationPath.removeLast()
        } else {
            navigationPath.append(AppRoutes.homeScreen)
        }
    }
}
