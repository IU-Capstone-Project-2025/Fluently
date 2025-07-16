//
//  WordOfTheDay.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI
import SwiftData

struct WordOfTheDay: View {
    @Environment(\.modelContext) var modelContext

    // MARK: - Properties
    var word: WordModel

    @State var inLibAlready = false

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
                .id(inLibAlready)
        }
        .onAppear() {
            inLibAlready = isWordInDictionary()
        }
    }

    func isWordInDictionary() -> Bool {
        let predicate = #Predicate<WordModel> { $0.id == word.id }
        let fetchDescriptor = FetchDescriptor<WordModel>(predicate: predicate)

        do {
            let results = try modelContext.fetch(fetchDescriptor)
            print(results)
            return !results.isEmpty
        } catch {
            print("Error checking word: \(error)")
            return false
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
            Image(systemName: !inLibAlready ? "plus.circle" : "checkmark.circle")
                .foregroundStyle(.blackText)
            Text( !inLibAlready ? "Add to collection" : "Already in collection")
                .foregroundStyle(.blackText)
                .font(.appFont.secondarySubheadline)
        }
        .onTapGesture {
            modelContext.insert(word)

            try? modelContext.save()
            withAnimation(.easeIn(duration: 0.3)) {
                inLibAlready = true
            }
        }
        .disabled(inLibAlready)
        .padding(.vertical, 4)
        .padding(.horizontal, 2)
        .background(
            RoundedRectangle(cornerRadius: Const.buttonCornerRadius)
                .fill(.grayFluently.opacity(0.4))
        )
    }
}
