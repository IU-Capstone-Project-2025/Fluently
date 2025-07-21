//
//  ProfileScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI
import SwiftData

protocol ProfileScreenPresenting: ObservableObject {

}

final class ProfileScreenPresenter: ProfileScreenPresenting {
    let router: ProfileScreenRouter
    let interactor: ProfileScreenInteractor

    @ObservedObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    @Published var preferences: PreferencesModel?

    @Published var goals: [String]
    @Published var dailyWord: Bool = true
    @Published var notifications: Bool = false
    @Published var notificationAt: Date = Date.now
    @Published var wordsPerDay: Int = 10
    @Published var goal: String = ""

    var modelContext: ModelContext?

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

        self.goals = []
    }

    // Navigation
    func navigateBack() {
        updatePrefs()
        router.navigateBack()
    }

    func setupPrefs(_ preferences: PreferencesModel?) {
        guard let preferences else {
            return
        }
        
        self.preferences = preferences

        dailyWord = preferences.dailyWord
        notifications = preferences.notifications
        notificationAt = preferences.notificationAt
        wordsPerDay = preferences.wordPerDay
        goal = preferences.goal

        getGoals()
    }

    func getGoals() {
        Task {
            do {
                goals = try await interactor.getGoals()
            } catch {
                print("Error while requesting goals: \(error)")
            }
        }
    }

    func getPrefs() {
        Task {
            do {

                let descriptor = FetchDescriptor<PreferencesModel>()
                let oldPrefs = try? modelContext?.fetch(descriptor)

                if let oldPrefs {
                    oldPrefs.forEach { pref in
                        modelContext?.delete(pref)
                    }
                }

                try modelContext?.save()

                preferences = try await interactor.getPreferences()
                if let prefs = preferences {
                    setupPrefs(prefs)
                    modelContext?.insert(prefs)

                    try modelContext?.save()
                }

            } catch {
                print("Error while fethcing preferences: \(error.localizedDescription)")
            }
        }
    }

    func updatePrefs() {
        guard let preferences else {
            return
        }

        print("saving")
        Task {
            do {
                preferences.dailyWord = dailyWord
                preferences.notificationAt = notificationAt
                preferences.notifications = notifications
                preferences.wordPerDay = wordsPerDay
                preferences.goal = goal

                try await interactor.api.updatePreferences(preferences)
            } catch {
                print("Error while fethcing preferences: \(error.localizedDescription)")
            }
        }
    }

    // Account
    func signOut() {
        authViewModel.signOut()
        account.isLoggedIn = false

        router.navigateBack()
    }
}
