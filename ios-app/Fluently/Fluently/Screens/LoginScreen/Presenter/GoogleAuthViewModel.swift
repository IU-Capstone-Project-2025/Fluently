//
//  GoogleAuthViewModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.06.2025.
//

import GoogleSignIn
import SwiftUI

class GoogleAuthViewModel: ObservableObject {
    // MARK: - Properties
    private weak var account: AccountData?

    @Published var isSignedIn = false
    @Published var userEmail: String?
    @Published var userName: String?

    // Sign in implementation
    func handleSignInButton(rootViewController: UIViewController?){
        if let rootViewController = rootViewController {
            GIDSignIn.sharedInstance.signIn(
                withPresenting: rootViewController) { signInResult, error in
                    guard let result = signInResult else {
                        if let error {
                            print("Error while signin: \(error)")
                            return
                        }
                        print("Error while signin: Unknown")
                        return
                    }
                    self.updateAccount(
                        name: result.user.profile?.name,
                        familyName: result.user.profile?.familyName,
                        email: result.user.profile?.email,
                        image: result.user.profile?.imageURL(withDimension: 100)?.absoluteString
                    )
                    print(result.user.idToken as Any)
                }
        }
    }

    // MARK: - Account interactions
    func setup(account: AccountData) {
        self.account = account
    }
    
    func updateAccount(
        name: String?,
        familyName: String?,
        email: String?,
        image: String?
    ) {
        self.account?.name = name ?? ""
        self.account?.familyName = familyName ?? ""
        self.account?.mail = email ?? ""
        self.account?.image = image

        self.isSignedIn = true
    }

    func signOut() {
        GIDSignIn.sharedInstance.signOut()
        isSignedIn = false
        userEmail = nil
        userName = nil
        print("User signed out")
    }
}

// MARK: - Errors
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
