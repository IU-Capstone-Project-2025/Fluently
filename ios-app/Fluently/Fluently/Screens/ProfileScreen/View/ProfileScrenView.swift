//
//  ProfileScrenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI
import SwiftData

struct ProfileScrenView: View {
    @ObservedObject var presenter: ProfileScreenPresenter

    @Query var prefs: [PreferencesModel]

    @Environment(\.modelContext) var modelContext

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
        .onAppear {
            presenter.modelContext = modelContext
            presenter.getPrefs()
//            presenter.setupPrefs(prefs.first)
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
            }
            Text(presenter.account.mail ?? "")
                .foregroundStyle(.orangeSecondary)
                .font(.appFont.secondaryCaption)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {
            userPreferences
        }
        .modifier(SheetViewModifier())
    }

    /// Sign out
    var signOutButton: some View {
        Button {
            presenter.signOut()
        } label: {
            Text("sign out")
                .frame(
                    alignment: .center
                )
                .frame(maxWidth: .infinity)
                .padding(.vertical, 10)
                .font(.appFont.title2)
                .foregroundStyle(.pink)
                .glass(
                    cornerRadius: 20,
                    fill: .pink
                )
        }
    }

    @ViewBuilder
    var userPreferences: some View {
        if let prefs = presenter.preferences {
            ScrollView {
                VStack(spacing: 12) {
                    cefrLevel(prefs.cefrLevel)
                    goal($presenter.goal)
                    settings(
                        dailyWord: $presenter.dailyWord,
                        notifications: $presenter.notifications
                    )
                    if presenter.notifications {
                        withAnimation(.easeInOut(duration: 0.3)) {
                            datePicker(date: $presenter.notificationAt)
                        }
                    }
                    numberOfWords(wordsNumber: $presenter.wordsPerDay)
                    Spacer(
                        minLength: 80
                    )

                    signOutButton
                        .frame(
                            maxHeight: .infinity,
                            alignment: .bottom
                        )
                        .safeAreaPadding()
                }
                .padding(.horizontal)
            }
            .scrollDismissesKeyboard(.interactively)
            .scrollIndicators(.hidden)
        } else {
            ZStack {
                VStack {
                    Text("No preferences loaded")
                        .font(.appFont.title)
                        .foregroundStyle(.grayFluently)
                    Text("Try to restart app, or check network connection")
                        .multilineTextAlignment(.center)
                        .font(.appFont.title3)
                        .foregroundStyle(.grayFluently)
                }
                .frame(
                    maxHeight: .infinity,
                    alignment: .center
                )

                signOutButton
                    .frame(
                        maxHeight: .infinity,
                        alignment: .bottom
                    )
                    .safeAreaPadding()
            }
            .padding()
        }
    }

    func cefrLevel(_ cefrLevel: String) -> some View {
        HStack {
            Text("Your CEFR Level:")
                .font(.appFont.title2)
                .frame(alignment: .leading)
            Spacer()
            Text(cefrLevel)
                .font(.appFont.largeTitle)
                .frame(alignment: .trailing)
        }
        .padding()
        .frame(maxWidth: .infinity)
        .foregroundStyle(.orangePrimary)
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }

    func goal(_ goal: Binding<String>) -> some View {
        HStack {
            Text("Your learning goal:")
                .font(.appFont.title2)
                .frame(alignment: .leading)
            Spacer()
            Picker(
                "Goal",
                selection: goal
            ) {
                ForEach(presenter.goals, id: \.hashValue) { topic in
                    Text(topic)
                        .tag(topic)
                }
            }
            .pickerStyle(.menu)
//            Text(goal)
//                .font(.appFont.title)
//                .frame(alignment: .trailing)
        }
        .padding()
        .frame(maxWidth: .infinity)
        .foregroundStyle(.orangePrimary)
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }

    func settings(
        dailyWord: Binding<Bool>,
        notifications: Binding<Bool>
    ) -> some View {
        VStack {
            Toggle(isOn: dailyWord) {
                Text("Daily words")
                    .font(.appFont.title2)
            }
            Toggle(isOn: notifications) {
                Text("Notifications")
                    .font(.appFont.title2)
            }
        }
        .padding()
        .foregroundStyle(.orangePrimary)
        .toggleStyle(.switch)
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }

    func datePicker(date: Binding<Date>) -> some View {
        DatePicker(
            "Notification time",
            selection: date,
            displayedComponents: [.hourAndMinute]
        )
        .font(.appFont.title2)
        .foregroundStyle(.orangePrimary)
        .datePickerStyle(.compact)
        .padding()
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }

    func numberOfWords(wordsNumber: Binding<Int>) -> some View {
        HStack {
            Text("Words per lesson")
                .foregroundStyle(.orangePrimary)
                .font(.appFont.title2)
            Spacer()

            TextField("words", value: wordsNumber, formatter: NumberFormatter())
                .keyboardType(.numberPad)
                .frame(width: 60)
                .padding(.vertical, 8)
                .padding(.horizontal, 3)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(.orangeSecondary)
                )
                .font(.body)
        }
        .frame(maxWidth: .infinity)
        .padding()
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }
}

// MARK: - Preview Provider
struct ProfileScreenPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        let router = AppRouter()
        let account = AccountData()
        let authVM = GoogleAuthViewModel()

        let profileView: ProfileScrenView

        init () {
            self.profileView =  ProfileScreenBuilder.build(
                router: self.router,
                account: self.account,
                authViewModel: self.authVM
            )

            self.profileView.presenter.preferences = PreferencesModel.generate()
        }

        var body: some View {
            profileView
            .environmentObject(account)
        }
    }
}
