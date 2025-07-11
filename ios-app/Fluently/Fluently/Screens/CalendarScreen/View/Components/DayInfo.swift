//
//  DayInfo.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI

struct DayInfo: View {
    var selectedDate: Date

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
            Spacer()
        }
        .frame(maxWidth: .infinity)
        .padding()
        .glass(
            cornerRadius: 20,
            fill: .orangePrimary
        )
    }
}
