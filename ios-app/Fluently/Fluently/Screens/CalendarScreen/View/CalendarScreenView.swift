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
