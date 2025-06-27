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

    func getPath() -> NavigationPath {
        return navigationPath
    }

    func navigate(to destination: any Hashable) {
        navigationPath.append(destination)
    }

    func popToRoot() {
        navigationPath.removeLast(navigationPath.count)
    }

    func pop() {
        navigationPath.removeLast()
    }
}
