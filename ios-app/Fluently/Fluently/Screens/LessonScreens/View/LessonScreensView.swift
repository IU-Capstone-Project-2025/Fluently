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
                Text(presenter.currentEx.exercizeType.rawValue)
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
            switch presenter.currentEx.exercizeType {
                case .chooseTranslationEngRuss:
                    Text(presenter.currentEx.exercizeType.rawValue)
                case .chooseTranslationRussEng:
                    Text(presenter.currentEx.exercizeType.rawValue)
                case .pickOptions:
                    let pickOptionEx = presenter.currentEx as! PickOptionsExs
                    PickOptionsView(
                        sentence: pickOptionEx.sentence,
                        answers: pickOptionEx.options) { selectedAnswer in
                            presenter.answer(selectedAnswer)
                        }
                        .id(presenter.currentEx.exercizeID)
                case .recordPronounce:
                    Text(presenter.currentEx.exercizeType.rawValue)
                case .wordCard:
                    Text(presenter.currentEx.exercizeType.rawValue)
                case .numberOfWords:
                    Text(presenter.currentEx.exercizeType.rawValue)
            }
        }
        .modifier(SheetViewModifier())
    }
}
