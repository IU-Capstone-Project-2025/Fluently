//
//  ProfileScrennView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct ProfileScrennView: View {
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var account: AccountData
    @ObservedObject var authViewModel: GoogleAuthViewModel

    @Binding var navigationPath: NavigationPath

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

    var logoutButton: some View {
        Button {
            authViewModel.signOut()
            account.isLoggined = false
            router.popToRoot()
        } label: {
            Text("sign out")
                .foregroundStyle(.red)
                .font(.title)
        }
    }

}
