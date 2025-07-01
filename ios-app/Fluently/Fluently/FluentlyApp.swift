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

#if targetEnvironment(simulator)

    @State private var showLogin = false
    @State private var showLaunchScreen = false

#else

    @State private var showLogin = true
    @State private var showLaunchScreen = true

#endif

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
                            HomeScreenBuilder.build(router: router, acoount: account)
                        } else {
                            LoginScreenBuilder.build(
                                router: router,
                                acount: account,
                                authViewModel: authViewModel
                            )
                                .onOpenURL(perform: handleURL)
                        }
                    }
                }
                .navigationDestination(for: AppRoutes.self) { route in
                    switch route {
                        case .homeScreen:
                            HomeScreenBuilder.build(
                                router: router,
                                acoount: account
                            )
                        case .login:
                            LoginScreenBuilder.build(
                                router: router,
                                acount: account,
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
                                router: router
                            )
                    }
                }
            }
            .environmentObject(account)
            .environmentObject(router)
            .modelContainer(for: Word.self)
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
        Task {
            do {
                let response = try await apiService.authGoogle(token!)
                print(response.accessToken)
                print(response.refreshToken)
                print(response.expiresIn)
                print(response.tokenType)
            } catch {
                print("error: \(error)")
            }
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
