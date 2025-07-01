//
//  LoginScreenBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

enum LoginScreenBuilder {

    static func build(
        router: AppRouter,
        acount: AccountData,
        authViewModel: GoogleAuthViewModel
    ) -> LoginView {
        let router = LoginRouter(router: router)
        let presenter = LoginPresenter (
            router: router,
            account: acount,
            authViewModel: authViewModel
        )
        
        return LoginView(
            presenter: presenter
        )
    }
}
