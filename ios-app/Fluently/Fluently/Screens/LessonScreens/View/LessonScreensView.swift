//
//  LessonScreensView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct LessonScreensView: View {
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var account: AccountData

    @State var showExitAlert = false
    @ObservedObject var presenter: LessonsPresenter

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
                Text("Lesson:")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
                Text(presenter.currentEx.exerciseType.rawValue)
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle.bold())
            }
            Spacer()
        }
        .padding(Const.horizontalPadding)
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack {
            switch presenter.currentEx.exerciseType {
                case .chooseTranslationEngRuss:
                    let chooseWordEx = presenter.currentEx as! ChooseTranslationExs
                    ChooseTranslationView(
                        word: chooseWordEx.word,
                        answers: chooseWordEx.options
                    ) { selectedAnswer in
                            presenter.answer(selectedAnswer)
                        }
                        .id(presenter.currentEx.exerciseId)
                case .typeTranslationRussEng:
                    let typeTranslationEx = presenter.currentEx as! TypeTranslationExs
                    TypeTranslationView (
                        word: typeTranslationEx.word) { typedAnswer in
                            presenter.answer(typedAnswer)
                        }
                case .pickOptions:
                    let pickOptionEx = presenter.currentEx as! PickOptionsExs
                    PickOptionsView(
                        sentence: pickOptionEx.sentence,
                        answers: pickOptionEx.options
                    ) { selectedAnswer in
                            presenter.answer(selectedAnswer)
                        }
                        .id(presenter.currentEx.exerciseId)
                case .recordPronounce:
                    Text(presenter.currentEx.exerciseType.rawValue)
                case .wordCard:
                    let wordCard = presenter.currentEx as! WordCard
                    WordCardView(
                        word: wordCard,
                        onKnowTapped: {
                            presenter.nextExercize()
                        },
                        onLearnTapped: {
                            presenter.showLesson()
                        }
                    )
                    .id(presenter.currentEx.exerciseId)
                case .numberOfWords:
                    Text(presenter.currentEx.exerciseType.rawValue)
            }
        }
        .modifier(SheetViewModifier())
    }
}
