//
//  HomeScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation

enum HomeScreenBuilder{
    static func build(
        router: AppRouter,
        acoount: AccountData,
    ) -> HomeScreenView {

        var presenter = HomeScreenPresenter(
            router: router,
            account: acoount
        )

        return HomeScreenView (
            presenter: presenter
        )
    }
}
