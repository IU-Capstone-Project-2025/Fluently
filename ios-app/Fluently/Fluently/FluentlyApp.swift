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
    // MARK: - Key parts
    @StateObject private var account = AccountData()
    @StateObject private var authViewModel = GoogleAuthViewModel()
    @StateObject private var router = AppRouter()


    @State private var showLogin = false

    var body: some Scene {
        WindowGroup {
            NavigationStack(path: $router.navigationPath) {
                Group {
                    if account.isLoggedIn && !showLogin {
                        HomeScreenBuilder.build(router: router, acoount: account)
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
                            HomeScreenBuilder.build(router: router, acoount: account)
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
                        case .lesson:
                            LessonScreensView(
                                presenter: LessonsPresenter(
                                    router: router,
                                    words: WordCardGenerator.generateCards()
                                )
                            )
                    }
                }
            }
            .onChange(of: account.isLoggedIn) {
                print("account is: \(account.isLoggedIn)")
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
                    account.image = user.profile?.imageURL(withDimension: 100)?.absoluteString
                    account.isLoggedIn = true
                    showLogin = false

                    print(user.idToken)
                } else {
                    account.isLoggedIn = false
                    showLogin = true
                }
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
