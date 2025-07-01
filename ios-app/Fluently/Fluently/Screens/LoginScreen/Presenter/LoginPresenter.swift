//
//  LoginPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI
import UIKit

protocol LoginPresenting: ObservableObject {
    // Auth
    func authViaGoogle(rootVC vc: UIViewController?)

    // Navigation
    func navigateToHome()
}

final class LoginPresenter: LoginPresenting {
    // MARK: - Key Objects
    var router: LoginRouter

    @ObservedObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    init(
        router: LoginRouter,
        account: AccountData,
        authViewModel: GoogleAuthViewModel
    ) {
        self.router = router
        self.account = account
        self.authViewModel = authViewModel
    }

    func navigateToHome() {
        router.navigateToHome()
    }

    func authViaGoogle(rootVC vc: UIViewController?) {
        authViewModel.setup(account: account)

        authViewModel.handleSignInButton(
            rootViewController: vc
        )
    }
}
