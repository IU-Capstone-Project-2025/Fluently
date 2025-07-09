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
                infoGrid
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
        VStack (alignment: .center) {

        }
        .modifier(SheetViewModifier())
    }
}

struct DaysHeader: View {

    @Binding var currentDate: Date

    var weekDayName: String {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "E"
        return dateFormatter.string(from: currentDate)
    }

    var body: some View {
        HStack() {

        }
    }


}
