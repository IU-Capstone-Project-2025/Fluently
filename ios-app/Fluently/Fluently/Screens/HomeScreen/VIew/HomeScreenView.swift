//
//  HomeScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI
import SwiftData

struct HomeScreenView: View {
    // MARK: - Key Objects
    @StateObject var presenter: HomeScreenPresenter
    @Environment(\.modelContext) var modelContext

    var words: [WordModel] {
        let descriptor = FetchDescriptor<WordModel>(
            predicate: #Predicate {
                $0.isInLesson == false &&
                $0.isInLibrary == true
            }
        )
        return (try? modelContext.fetch(descriptor)) ?? []
    }

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
        .onAppear {
            presenter.modelContext = modelContext
            presenter.getDayWord()
            Task {
                do {
                    try await presenter.getLesson()
                    presenter.compare()
                    await presenter.checkForNilIDs()
                } catch {
                    print(error)
                }
            }
            presenter.compare()
        }
        .navigationBarBackButtonHidden()
        .modifier(BackgroundViewModifier())

        .fullScreenCover(item: $openedScreen) { screenType in
            switch screenType {
                case .notes:
                    presenter.buildNotesScreen()
                case .learned:
                    presenter.buildDictionaryScreen(isLearned: true)
                case .nonLearned:
                    presenter.buildDictionaryScreen(isLearned: false)
            }
        }
    }

    // MARK: - SubViews

    ///  Top Bar
    var topBar: some View {
        HStack {
            VStack (alignment: .leading) {
                Text("Home")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
//                Text(goal)
//                    .foregroundStyle(.whiteText)
//                    .font(.appFont.largeTitle.bold())
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
            WordOfTheDay(word: presenter.wordOfTheDay ?? WordModel.mockWord())
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
            ScreenCard(type: .notes) {
                openedScreen = .notes
            }
            ScreenCard(type: .learned, count: "\(words.filter { $0.isLearned == true}.count)") {
                openedScreen = .learned
            }
            ScreenCard(type: .nonLearned, count: "\(words.filter { $0.isLearned == false}.count)") {
                openedScreen = .nonLearned
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
                .id(presenter.lesson?.id)
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
