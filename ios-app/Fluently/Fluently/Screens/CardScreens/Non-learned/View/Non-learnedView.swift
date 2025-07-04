//
//  Non-learnedView.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

struct NonLearnedView: View {
    @EnvironmentObject var router: AppRouter

    @Environment(\.dismiss) var dismiss

    // MARK: - Constants
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
            .toolbar {
                ToolbarItem(placement: .topBarLeading) {
                    Button {
                        dismiss.callAsFunction()
                    } label: {
                        Image(systemName: "chevron.left")
                        .foregroundStyle(.whiteText)
                    }
                }
            }
        }
    }

    // MARK: - SubViews

    /// Top Bar
    var topBar: some View {
        VStack(alignment: .center) {
            Text("Non-Learned")
                .foregroundStyle(.whiteText)
                .font(.appFont.largeTitle.bold())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.horizontal, Const.horizontalPadding)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {
            
        }
        .modifier(SheetViewModifier())
    }
}

struct NonLearnedPreview: PreviewProvider {

    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        var body: some View {
            NonLearnedView()
        }
    }
}
