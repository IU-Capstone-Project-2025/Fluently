//
//  ProfileScrenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct ProfileScrenView: View {
    @ObservedObject var presenter: ProfileScreenPresenter

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
                    presenter.navigateBack()
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
                Text(presenter.account.name ??  "")
                    .foregroundStyle(.orangeSecondary)
                    .font(.appFont.secondaryCaption)
//                Text(account.familyName ??  "")
//                    .foregroundStyle(.orangeSecondary)
//                    .font(.appFont.secondaryCaption)
            }
            Text(presenter.account.mail ?? "")
                .foregroundStyle(.orangeSecondary)
                .font(.appFont.secondaryCaption)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {
            signOutButton
        }
        .modifier(SheetViewModifier())
    }

    /// Sign out
    var signOutButton: some View {
        Button {
            presenter.signOut()
        } label: {
            Text("sign out")
                .foregroundStyle(.red)
                .font(.title)
        }
    }

}
