//
//  DayInfo.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI

struct DayInfo: View {
    var selectedDate: Date
    var words: [WordModel]

    @State private var isLearned = true

    var dateLabel: String {
        let formatter = DateFormatter()
        formatter.dateFormat = "d MMMM"
        return formatter.string(from: selectedDate)
    }

    var yearFormatter: String {
        let formatter = DateFormatter()
        formatter.dateFormat = "YYYY"
        return formatter.string(from: selectedDate)
    }

    var body: some View {
        infoGrid
            .padding()
    }

    // MARK: - SubViews

    /// Main layer with info
    var infoGrid: some View {
        VStack(alignment: .center) {
            HStack {
                VStack {
                    Text(dateLabel)     /// label with selected day date
                        .foregroundStyle(.blackText)
                        .font(.appFont.largeTitle)
                        .frame(
                            maxWidth: .infinity,
                            alignment: .leading
                        )
                    Text(yearFormatter)  /// label with selected day year
                        .foregroundStyle(.blackText)
                        .font(.appFont.largeTitle)
                        .frame(
                            maxWidth: .infinity,
                            alignment: .leading
                        )
                }
                VStack {
                    Text("You've Learned")
                        .foregroundStyle(.blackText)
                        .font(.appFont.title2)
                        .frame(
                            maxWidth: .infinity,
                            alignment: .leading
                        )
                    /// number of leaned worsd at this day
                    Text("\(words.filter({ $0.isLearned == true}).count) words")
                        .foregroundStyle(.blackText)
                        .font(.appFont.title2)
                        .frame(
                            maxWidth: .infinity,
                            alignment: .leading
                        )
                }
            }
            /// selector of `leaned / non-leaner` filter
            Picker("Filter Words", selection: $isLearned) {
                Text("Learned").tag(true)
                    .foregroundStyle(.blackText)
                Text("Not Learned").tag(false)
                    .foregroundStyle(.blackText)
            }
            .pickerStyle(.segmented)

            learnedNonLearnedWords

            Spacer()
        }
        .frame(maxWidth: .infinity)
        .padding()
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }

    /// list of filtred words
    var learnedNonLearnedWords: some View {
        ScrollView {
            VStack(spacing: 10) {
                ForEach(words.filter({ $0.isLearned == isLearned}) , id: \.wordId) { word in
                    WordCardRow(word: word)
                }
            }
            .padding()
        }
        .scrollIndicators(.hidden)
    }
}
