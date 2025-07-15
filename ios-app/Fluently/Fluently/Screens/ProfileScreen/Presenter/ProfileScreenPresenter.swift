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
    let interactor: ProfileScreenInteractor

    @ObservedObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel
//#if targetEnvironment(simulator)
//    @Published var preferences: PreferencesModel = PreferencesModel.generate()
//#else
    @Published var preferences: PreferencesModel?
//#endif

    @State var dailyWord: Bool = true
    @State var notifications: Bool = false

    init(
        router: ProfileScreenRouter,
        interactor: ProfileScreenInteractor,
        account: AccountData,
        authViewModel: GoogleAuthViewModel
    ) {
        self.router = router
        self.interactor = interactor
        self.account = account
        self.authViewModel = authViewModel

        getPrefs()
    }

    // Navigation
    func navigateBack() {
        router.navigateBack()
    }

    func getPrefs() {
        Task {
            do {
                preferences = try await interactor.getPreferences()
            } catch {
                print("Error while fethcing preferences: \(error.localizedDescription)")
            }
        }
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
