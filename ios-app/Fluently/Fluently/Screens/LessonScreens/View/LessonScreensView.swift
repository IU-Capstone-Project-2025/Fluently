//
//  LessonScreensView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct LessonScreensView: View {
    // MARK: - Key objects
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var account: AccountData

    @State var showExitAlert = false
    @ObservedObject var presenter: LessonsPresenter

    // MARK: - View Constances
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
        .toolbar() {
            /// exit button
            ToolbarItem(placement: .topBarLeading) {
                Button {
                    showExitAlert = true
                } label: {
                    Image(systemName: "chevron.left")
                        .foregroundStyle(.whiteText)
                }
            }
        }
        .alert("Are you sure, that you want exit?", isPresented: $showExitAlert) {
            Button ("No", role: .cancel) {
                showExitAlert = false
            }
            Button ("Yes", role: .destructive) {
                router.pop()
            }
        }
        .navigationBarBackButtonHidden()
        .modifier(BackgroundViewModifier())
    }

    // MARK: - SubViews

    ///  Top Bar
    var topBar: some View {
        HStack {
            VStack (alignment: .leading) {
#if targetEnvironment(simulator)
                Text("Exercise:")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
                Text(presenter.currentEx.type.rawValue)
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
#else
                Text("Exercise: \(presenter.currentExNumber + 1)/\(presenter.words.count)")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
#endif
            }
            Spacer()
        }
        .padding(.horizontal, Const.horizontalPadding)
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack {
            Spacer()
                .frame(height: 80) // Hard code :(
            switch presenter.currentExType {
                case .chooseTranslationEngRuss: /// Choose correct translation
                    let chooseWordEx = presenter.currentEx.data as! ChooseTranslationEngRuss
                    ChooseTranslationView(
                        word: chooseWordEx.text,
                        answers: chooseWordEx.options
                    ) { selectedAnswer in
                            presenter.answer(selectedAnswer)
                        }
                    .id(presenter.currentExNumber)
                case .typeTranslationRussEng: /// Type correct translation
                    let typeTranslationEx = presenter.currentEx.data as! WriteFromTranslation
                    TypeTranslationView (word: typeTranslationEx.translation) { typedAnswer in
                            presenter.answer(typedAnswer)
                        }
                case .pickOptionSentence: /// Pick word, mathing by definition
                    let pickOptionEx = presenter.currentEx.data as! PickOptionSentence
                    PickOptionsView(
                        sentence: pickOptionEx.template,
                        answers: pickOptionEx.options
                    ) { selectedAnswer in
                            presenter.answer(selectedAnswer)
                        }
                    .id(presenter.currentExNumber)
                case .recordPronounce:
                    Text(presenter.currentEx.type.rawValue)
                case .wordCard:
                    let wordCard = presenter.words[presenter.currentExNumber]
                    WordCardView(
                        word: wordCard,
                        onKnowTapped: {
                            presenter.nextExercise()
                        },
                        onLearnTapped: {
                            presenter.showLesson()
                        }
                    )
                    .id(presenter.currentExNumber)
                case .numberOfWords:
                    Text(presenter.currentEx.type.rawValue)
            }
        }
        .modifier(SheetViewModifier())
    }
}
