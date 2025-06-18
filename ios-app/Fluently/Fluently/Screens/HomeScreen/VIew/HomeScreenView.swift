//
//  HomeScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

struct HomeScreenView: View {

    // MARK: - Properties
    @State var goal: String = "Traveling"
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
        .containerRelativeFrame([.horizontal, .vertical])
        .background(.orangePrimary)
    }

    // MARK: - SubViews

    ///  Top Bar
    var topBar: some View {
        HStack {
            VStack (alignment: .leading) {
                Text("Goal:")
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle)
                Text(goal)
                    .foregroundStyle(.whiteText)
                    .font(.appFont.largeTitle)
            }
            Spacer()
            AvatarImage(
                size: 100
            )
        }
        .padding(Const.horizontalPadding)
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack {
            WordOfTheDay(word: Word.mockWord())
            cards

            Spacer()

            LessonInfo(minutes: 10, seconds: 20)
            startButton

            Spacer()
        }
        .padding(.top, Const.gridInfoVerticalPadding)
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .background(
            UnevenRoundedRectangle(
                topLeadingRadius: Const.sheetCornerRadius,
                topTrailingRadius: Const.sheetCornerRadius,
            )
            .fill(
                .whiteBackground
            )
            .ignoresSafeArea(.all)
        )
    }

    ///  List of the cards
    var cards: some View {
        HStack(spacing: 12) {
            ForEach(ScreenCard.CardType.allCases, id: \.self) { type in
                ScreenCard(type: type) {
                    print(type.rawValue)
                }
            }
        }
        .padding(.horizontal, Const.horizontalPadding)
    }

    /// button to start lesson
    var startButton: some View {
        Text("Start")
            .foregroundStyle(.whiteText)
            .font(.appFont.title2)
            .padding(.vertical, 6)
            .frame(maxWidth: .infinity)
            .background(
                RoundedRectangle(cornerRadius: 50)
                    .fill(.blackFluently)
            )
            .padding(.horizontal, Const.horizontalPadding * 3)
    }
}

struct NavigationBar: View {

    var body: some View {
         Text("bottom bar")
    }
}
