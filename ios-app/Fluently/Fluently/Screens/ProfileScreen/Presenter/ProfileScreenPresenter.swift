//
//  ProfileScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

protocol ProfileScreenPresenting: ObservableObject {

}

final class ProfileScreenPresenter: ProfileScreenPresenting {
    let router: ProfileScreenRouter

    @ObservedObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    init(
        router: ProfileScreenRouter,
        account: AccountData,
        authViewModel: GoogleAuthViewModel
    ) {
        self.router = router
        self.account = account
        self.authViewModel = authViewModel
    }

    // Navigation
    func navigateBack() {
        router.navigateBack()
    }

    // Account
    func signOut() {
        authViewModel.signOut()
        account.isLoggedIn = false

        // TODO: - Think more abour this implementation
        router.popToRoot()
        router.navigateToLogin()
    }
}
