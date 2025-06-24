//
//  FluentlyApp.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import SwiftUI
import GoogleSignIn

@main
struct FluentlyApp: App {
    @StateObject private var account = AccountData()
    @StateObject private var authViewModel = GoogleAuthViewModel()
    @StateObject private var router = AppRouter()


    @State private var showLogin = false

    var body: some Scene {
        WindowGroup {
            NavigationStack(path: $router.navigationPath) {
                Group {
                    if account.isLoggined && !showLogin {
                        HomeScreenView()
                            .onDisappear {
                                showLogin = false
                            }
                    } else {
                        LoginView(
                            authViewModel: authViewModel,
                            navigationPath: $router.navigationPath
                        )
                            .onOpenURL(perform: handleURL)
                            .onAppear() {
                                attemptRestoreLogin()
                            }
                    }
                }
                .navigationDestination(for: AppRoutes.self) { route in
                    switch route {
                        case .homeScreen:
                            HomeScreenView()
                        case .login:
                            LoginView (
                                authViewModel: authViewModel,
                                navigationPath: $router.navigationPath
                            )
                        case .profile:
                            ProfileScrennView(
                                authViewModel: authViewModel,
                                navigationPath: $router.navigationPath
                            )
                    }
                }
            }
            .onChange(of: account.isLoggined) {
                print("account is: \(account.isLoggined)")
            }
            .environmentObject(account)
            .environmentObject(router)
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
                    account.isLoggined = true
                    showLogin = false
                } else {
                    account.isLoggined = false
                    showLogin = true
                }
            }
        }
    }
}


enum AppRoutes: Hashable {
    case homeScreen
    case login
    case profile
}
