//
//  FluentlyApp.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import SwiftUI
import GoogleSignIn
import SwiftData

@main
struct FluentlyApp: App {
    // MARK: - Key parts
    @StateObject private var account = AccountData()
    @StateObject private var authViewModel = GoogleAuthViewModel()
    @StateObject private var router = AppRouter()

    private var apiService = APIService()

    @Environment(\.modelContext) var modelContext

#if targetEnvironment(simulator)

    @State private var showLogin = false
    @State private var showLaunchScreen = false

#else

    @State private var showLogin = true
    @State private var showLaunchScreen = true

#endif

    var dayWordContainer: ModelContainer {
        do {
            let container = try ModelContainer(
                for: DayWord.self,
            )

            return container
        } catch {
            fatalError("Container for daily word is not created")
        }
    }

    init () {
        UIView.appearance().overrideUserInterfaceStyle = .light /// brutу force light theme
    }

    var body: some Scene {
        WindowGroup {
            NavigationStack(path: $router.navigationPath) {
                Group {
                    if showLaunchScreen {
                        LaunchScreenView(isActive: $showLaunchScreen)
                            .onAppear() {
                                attemptRestoreLogin()
                            }
                    } else {
                        if !showLogin {
                            MainView()
                        } else {
                            LoginScreenBuilder.build(
                                router: router,
                                account: account,
                                authViewModel: authViewModel
                            )
                                .onOpenURL(perform: handleURL)
                        }
                    }
                }
                .navigationDestination(for: AppRoutes.self) { route in
                    switch route {
                        case .homeScreen:
                            MainView()
                        case .login:
                            LoginScreenBuilder.build(
                                router: router,
                                account: account,
                                authViewModel: authViewModel
                            )
                        case .profile:
                            ProfileScreenBuilder.build(
                                router: router,
                                account: account,
                                authViewModel: authViewModel
                            )
                        case .lesson:
                            LessonScreenBuilder.build(
                                router: router,
                            )
                    }
                }
            }
            .environmentObject(account)
            .environmentObject(router)
            .modelContainer(
                for: [
                    CardsModel.self,
                    WordModel.self,
                    DayWord.self,
                    PreferencesModel.self
                ]
            )
        }
    }

    private func handleURL(_ url: URL) {
        GIDSignIn.sharedInstance.handle(url)
    }

    private func attemptRestoreLogin() {
        GIDSignIn.sharedInstance.restorePreviousSignIn { user, error in
            DispatchQueue.main.async {
                if let user = user {
                    account.name = user.profile?.name
                    account.familyName = user.profile?.familyName
                    account.mail = user.profile?.email
                    account.image = user.profile?.imageURL(withDimension: 100)?.absoluteString
                    account.isLoggedIn = true
                    showLogin = false

                    /// Api token request
                    requestAccessTokens(token: user.idToken?.tokenString)
                } else {
                    account.isLoggedIn = false
                    showLogin = true
                }
            }
        }
    }

    private func requestAccessTokens(token: String?) {
        guard let token = token else {
            fatalError("The token is empty")
        }

        Task {
            do {
                let response = try await apiService.authGoogle(token)
                do {
                    try KeyChainManager.shared.saveToken(response)
                } catch {
                    print("token saving error: \(error)")
                }
                getPreferences()
            } catch {
                print("response receiving error: \(error)")
            }
        }
    }

    private func getPreferences() {
        let descriptor = FetchDescriptor<PreferencesModel>()
        do {
            let prefs: [PreferencesModel] = try modelContext.fetch(descriptor)
            prefs.forEach { pref in
                modelContext.delete(pref)
            }
            try modelContext.save()

            Task {
                let newPreferences = try? await apiService.getPreferences()
                if let newPreferences {
                    modelContext.insert(newPreferences)
                }
                try modelContext.save()
            }
        } catch {
            print("Error while getting preferences: \(error)")
        }
    }
}

// MARK: - Routes
enum AppRoutes: Hashable {
    case homeScreen
    case login
    case profile
    case lesson
}
