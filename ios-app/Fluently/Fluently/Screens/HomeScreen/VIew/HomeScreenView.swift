//
//  HomeScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

struct HomeScreenView: View {
    // MARK: - Key Objects
    @ObservedObject var presenter: HomeScreenPresenter

    // MARK: - Properties
    @State var goal: String = "Traveling"
    
    @State var openedScreen: ScreenCard.CardType?

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
        .task {
            do {
                // TODO: fix this logic
                try await presenter.getLesson()
            } catch {
                print(error)
            }
        }

        .fullScreenCover(item: $openedScreen) { screenType in
            switch screenType {
                case .notes:
                    presenter.buildNotesScreen()
                case .learned:
                    presenter.buildDictionaryScreen()
                case .nonLearned:
                    NonLearnedView()
            }
        }
    }

    // MARK: - SubViews

    ///  Top Bar
    var topBar: some View {
        HStack {
            VStack (alignment: .leading) {
                Text("Goal:")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
                Text(goal)
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
            }
            Spacer()
            AvatarImage(
                size: 100,
                onTap: {
                    presenter.navigatoToProfile()
                }
            )
        }
        .padding(Const.horizontalPadding)
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack {
            WordOfTheDay(word: WordModel.mockWord())
            cards

            Spacer()

            LessonInfo(minutes: 10, seconds: 20)
            startButton

            Spacer()
        }
        .modifier(SheetViewModifier())
    }

    ///  List of the cards
    var cards: some View {
        HStack(spacing: 12) {
            ForEach(ScreenCard.CardType.allCases, id: \.self) { type in
                ScreenCard(type: type) {
                    openedScreen = type
                }
                .animation(.easeInOut, value: openedScreen)
            }
        }
        .padding(.horizontal, Const.horizontalPadding)
    }

    /// button to start lesson
    var startButton: some View {
        Button {
            presenter.navigatoToLesson()
        } label: {
            Text(presenter.lesson == nil ? "loading" : "Start")
                .foregroundStyle(.whiteText)
                .font(.appFont.title2.bold())
                .padding(.vertical, 6)
                .frame(maxWidth: .infinity)
                .background(
                    RoundedRectangle(cornerRadius: 50)
                        .fill(.blackFluently)
                )
                .padding(.horizontal, Const.horizontalPadding * 3)
        }
        .disabled(presenter.lesson == nil)
    }
}


struct HomeScreenPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        @StateObject var router = AppRouter()
        @StateObject var account = AccountData()

        var body: some View {
            HomeScreenBuilder.build(
                router: router,
                acoount: account
            )
            .environmentObject(account)
        }
    }
}
