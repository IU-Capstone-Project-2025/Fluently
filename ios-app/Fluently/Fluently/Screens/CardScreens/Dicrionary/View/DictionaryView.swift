//
//  DictionaryView.swift
//  Fluently
//
//  Created by Савва Пономарев on 30.06.2025.
//


import SwiftUI

struct DictionaryView: View {
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
            Text("Dictionary")
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

struct DictionaryPreview: PreviewProvider {

    static var previews: some View {
        DictionaryPreviewWrapper()
    }

    struct DictionaryPreviewWrapper: View {
        @State private var path = NavigationPath()

        var body: some View {
            DictionaryView()
        }
    }
}
