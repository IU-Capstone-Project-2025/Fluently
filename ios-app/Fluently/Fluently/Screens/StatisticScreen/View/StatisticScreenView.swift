//
//  StatisticScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import Foundation
import SwiftUI

struct StatisticScreenView: View {

    @ObservedObject var presenter: StatisticScreenPresenter

    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)
    }

    var body: some View {
        NavigationStack {
            VStack {
                topBar
                ZStack {
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
            Text("Statistic")
                .foregroundStyle(.whiteText)
                .font(.appFont.largeTitle.bold())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding( Const.horizontalPadding)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {}
        .modifier(SheetViewModifier())
    }

    var infoLayer: some View {
        VStack {
            RangeHeader(selectedRange: $presenter.selectedRange)
            StatisticInfo(range: presenter.selectedRange)
        }
    }
}


struct StatisticScreenPreview: PreviewProvider {

    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        var body: some View {
            StatisticScreenBuilder.build()
        }
    }
}
