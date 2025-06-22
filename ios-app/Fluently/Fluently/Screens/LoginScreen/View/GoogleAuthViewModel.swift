//
//  GoogleAuthViewModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.06.2025.
//

import GoogleSignIn
import SwiftUI

class GoogleAuthViewModel: ObservableObject {
    @Published var isSignedIn = false
    @Published var userEmail: String?
    @Published var userName: String?

//    func handleSignInButton(rootViewController: UIViewController?){
//        if let rootViewController = rootViewController {
//            GIDSignIn.sharedInstance.signIn(
//                withPresenting: rootViewController) { signInResult, error in
//                    guard let result = signInResult else {
//                        return
//                    }
//                }
//        }
//    }

    private weak var account: AccountData?

    func setup(account: AccountData) {
        self.account = account
    }

    func handleSignInButton(rootViewController: UIViewController?) {
        guard let rootViewController = rootViewController else { return }

        GIDSignIn.sharedInstance.signIn(withPresenting: rootViewController) { [weak self] signInResult, error in
            DispatchQueue.main.async {
                guard let result = signInResult else {
                    print("Error signing in: \(error?.localizedDescription ?? "Unknown error")")
                    return
                }

                // Update account data
                self?.account?.name = result.user.profile?.name
                self?.account?.familyName = result.user.profile?.familyName
                self?.account?.mail = result.user.profile?.email
                self?.account?.isLoggined = true

                // Update local view model state
                self?.isSignedIn = true
                self?.userEmail = result.user.profile?.email
                self?.userName = result.user.profile?.name
            }
        }
    }

    func signOut() {
        GIDSignIn.sharedInstance.signOut()
        isSignedIn = false
        userEmail = nil
        userName = nil
        print("User signed out")
    }
}

enum AuthError: Error, LocalizedError {
    case userCancelled
    case unknown

    var errorDescription: String? {
        switch self {
        case .userCancelled: return "Sign in was cancelled"
        case .unknown: return "An unknown error occurred"
        }
    }
}
