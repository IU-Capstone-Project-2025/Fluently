//
//  CalendarScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import Foundation
import SwiftUI

struct CalendarScreenView: View {
    @ObservedObject var presenter: CalendarScreenPresenter

    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)
    }

    var body: some View {
        NavigationStack {
            VStack {
                topBar
                ZStack (alignment: .top) {
                    infoGrid
                    infoLayer
                }
            }
            .navigationBarBackButtonHidden()
            .modifier(BackgroundViewModifier())
        }
    }

    // MARK: - SubViews

    /// Top Bar
    var topBar: some View {
        VStack(alignment: .center) {
            Text("Calendar")
                .foregroundStyle(.whiteText)
                .font(.appFont.largeTitle.bold())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(Const.horizontalPadding)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {}
        .modifier(SheetViewModifier())
    }

    var infoLayer: some View {
        VStack {
            DaysHeader(selectedDate: $presenter.selectedDate)
                .glass(
                    cornerRadius: 0,
                    fill: .orangePrimary
                )
            DayInfo(selectedDate: presenter.selectedDate)
        }
    }
}

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



// MARK: - Preview Provider
struct CalendarScreenPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        let presenter = CalendarScreenPresenter()

        var body: some View {
            CalendarScreenView(presenter: presenter)
        }
    }
}
