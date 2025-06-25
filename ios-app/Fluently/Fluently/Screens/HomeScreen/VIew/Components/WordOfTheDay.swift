//
//  WordOfTheDay.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//


import SwiftUI

struct WordOfTheDay: View {
    // MARK: - Properties
    @State var word: Word

    // MARK: - Constants
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(40)
        static let verticalPadding = CGFloat(20)

        // Corner Radius
        static let cardCornerRadius = CGFloat(20)
        static let buttonCornerRadius = CGFloat(50)
    }

    // MARK: - Main Body
    var body: some View {
        VStack(alignment: .center, spacing: 4) {
            Text("Word of the day")
                .font(.appFont.title3.bold())

            wordCard
            addCard
        }
    }

    // MARK: - Subviews

    /// Card displaying the word of the day and its translation
    private var wordCard: some View {
        VStack {
            Text(word.word)
                .font(.appFont.title.bold())
                .foregroundStyle(.whiteText)

            Text(word.translation)
                .font(.appFont.callout)
                .foregroundStyle(.whiteBackground.secondary)
        }
        .padding(.vertical, Const.verticalPadding)
        .padding(.horizontal, Const.horizontalPadding)
        .background(
            RoundedRectangle(cornerRadius: Const.cardCornerRadius)
                .fill(.blackFluently)
        )
    }

    /// Button for adding the word to collection
    private var addCard: some View {
        HStack(spacing: 3) {
            Image(systemName: "plus.circle")
                .foregroundStyle(.blackText)
            Text("Add to collection")
                .foregroundStyle(.blackText)
                .font(.appFont.secondarySubheadline)
        }
        .padding(.vertical, 4)
        .padding(.horizontal, 2)
        .background(
            RoundedRectangle(cornerRadius: Const.buttonCornerRadius)
                .fill(.grayFluently.opacity(0.4))
        )
    }
}
