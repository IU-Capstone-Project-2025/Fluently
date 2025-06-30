//
//  ProfileScrennView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct ProfileScrennView: View {
    // MARK: - Key Objects
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    // MARK: - Properties
    @Binding var navigationPath: NavigationPath

    // MARK: - View Constances
    private enum Const {
        static var avatarSize = CGFloat(120)
    }

    var body: some View {
        VStack {
            topBar

            Spacer()
            infoGrid
        }
        .navigationBarBackButtonHidden()
        .modifier(BackgroundViewModifier())
        .toolbar {
            ToolbarItem(placement: .topBarLeading) {
                Button {
                    router.pop()
                } label: {
                    Image(systemName: "chevron.left")
                        .foregroundStyle(.whiteText)
                }
            }
        }
    }

    // MARK: - Subviews

    var topBar: some View {
        VStack {
            AvatarImage(size: Const.avatarSize)
            HStack {
                Text(account.name ??  "")
                    .foregroundStyle(.orangeSecondary)
                    .font(.appFont.secondaryCaption)
//                Text(account.familyName ??  "")
//                    .foregroundStyle(.orangeSecondary)
//                    .font(.appFont.secondaryCaption)
            }
            Text(account.mail ?? "")
                .foregroundStyle(.orangeSecondary)
                .font(.appFont.secondaryCaption)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {
            logoutButton
        }
        .modifier(SheetViewModifier())
    }

    /// Log out
    var logoutButton: some View {
        Button {
            authViewModel.signOut()
            account.isLoggedIn = false
            
            // TODO: - Think more abour this implementation
            router.popToRoot()
            router.navigate(to: AppRoutes.login)
        } label: {
            Text("sign out")
                .foregroundStyle(.red)
                .font(.title)
        }
    }

}
