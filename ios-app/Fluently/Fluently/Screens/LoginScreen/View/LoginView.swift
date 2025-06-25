//
//  LoginView.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.06.2025.
//

import Foundation
import SwiftUI
import UIKit

import GoogleSignInSwift
import GoogleSignIn

struct LoginView: View {
    // MARK: - Key Objects
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    @Binding var navigationPath: NavigationPath

    // MARK: - Properties
    let name: String = "Fluently"

    // MARK: - Constants
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)

        // Corner Radiuses
        static let sheetCornerRadius = CGFloat(20)
        static let gridInfoVerticalPadding = CGFloat(20)
    }

    var body: some View {
        VStack {
            topBar
            infoGrid
        }
        .navigationBarBackButtonHidden()
        .modifier(BackgroundViewModifier())
        .onReceive(authViewModel.$isSignedIn) { isSignedIn in
            if isSignedIn {
                navigationPath.append(AppRoutes.homeScreen)
            }
        }
    }

    // MARK: - SubViews

    /// Top Bar
    var topBar: some View {
        HStack {
            VStack (alignment: .center) {
                Text("Welcome to")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
                Text(name)
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
            }
        }
        .padding(Const.horizontalPadding)
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {
            googleSignInButton
        }
        .modifier(SheetViewModifier())
    }

    // MARK: - Google sign in implementation
    var googleSignInButton: some View {
        GoogleSignInButton(viewModel: GoogleSignInButtonViewModel(
            scheme: .dark,
            style: .wide,
            state: .normal
        )) {
            authViewModel.setup(account: account)
            
            authViewModel.handleSignInButton(
                rootViewController: getRootViewController()
            )
        }
        .padding(.horizontal, Const.horizontalPadding)
    }

    func getRootViewController() -> UIViewController? {
        guard let scene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
              let rootViewController = scene.windows.first?.rootViewController else {
            return nil
        }
        return getVisibleViewController(from: rootViewController)
    }

    private func getVisibleViewController(from vc: UIViewController) -> UIViewController {
        if let nav = vc as? UINavigationController {
            return getVisibleViewController(from: nav.visibleViewController!)
        }
        if let tab = vc as? UITabBarController {
            return getVisibleViewController(from: tab.selectedViewController!)
        }
        if let presented = vc.presentedViewController {
            return getVisibleViewController(from: presented)
        }
        return vc
    }
}

