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
            Text(dateLabel)
                .foregroundStyle(.blackText)
                .font(.appFont.largeTitle)
                .frame(
                    maxWidth: .infinity,
                    alignment: .leading
                )
            Text(yearFormatter)
                .foregroundStyle(.blackText)
                .font(.appFont.largeTitle)
                .frame(
                    maxWidth: .infinity,
                    alignment: .leading
                )
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

    var learnedNonLearnedWords: some View {
        ScrollView {
//            Text( isLearned ? "Learned" : "Non-Learned")
//                .font(.appFont.title)
//                .foregroundStyle(.blackText)
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
