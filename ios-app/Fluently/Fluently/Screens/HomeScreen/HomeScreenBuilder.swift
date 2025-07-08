//
//  HomeScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation

// MARK: - View Builder
enum HomeScreenBuilder{
    static func build(
        router: AppRouter,
        acoount: AccountData,
    ) -> HomeScreenView {
        let router = HomeScreenRouter(router: router)
        let interactor = HomeScreenInteractor()
        let presenter = HomeScreenPresenter(
            router: router,
            interactor: interactor,
            account: acoount
        )

        return HomeScreenView (
            presenter: presenter
        )
    }
}

