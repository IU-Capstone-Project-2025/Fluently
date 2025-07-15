//
//  ProfileScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

enum ProfileScreenBuilder {

    static func build (
        router: AppRouter,
        account: AccountData,
        authViewModel: GoogleAuthViewModel
    ) -> ProfileScrenView {
        let router = ProfileScreenRouter(router: router)
        let interactor = ProfileScreenInteractor()

        let presenter = ProfileScreenPresenter(
            router: router,
            interactor: interactor,
            account: account,
            authViewModel: authViewModel
        )

        return ProfileScrenView(
            presenter: presenter
        )
    }
}
