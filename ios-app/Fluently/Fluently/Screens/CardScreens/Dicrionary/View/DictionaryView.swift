//
//  DictionaryView.swift
//  Fluently
//
//  Created by Савва Пономарев on 30.06.2025.
//

import SwiftUI
import SwiftData

struct DictionaryView: View {
    @Environment(\.dismiss) var dismiss

    @Environment(\.modelContext) var modelContext

    @StateObject private var presenter: DictionaryScreenPresenter

    var isLearned: Bool

    init(isLearned: Bool) {
        self.isLearned = isLearned
        _presenter = StateObject(wrappedValue: DictionaryScreenPresenter(
            isLearned: isLearned
        ))
    }

    @State var prefix: String = ""

    // MARK: - Constants
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)
    }

    var body: some View {
        NavigationStack {
            GeometryReader { _ in
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
            .onAppear{
                presenter.setModelContext(modelContext)
                presenter.filter(prefix: "")
            }
        }
        .ignoresSafeArea(.keyboard)
    }

    // MARK: - SubViews

    /// Top Bar
    var topBar: some View {
        VStack(alignment: .center) {
            Text( isLearned ? "Learned" : "Non-Learned")
                .foregroundStyle(.whiteText)
                .font(.appFont.largeTitle.bold())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.horizontal, Const.horizontalPadding)
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack {
            SearchBar(text: $prefix)
                .padding(.bottom, 12)
                .padding(.horizontal)
            ScrollView {
                LazyVStack (alignment: .center, spacing: 12) {
                    ForEach(presenter.filteredWords, id: \.wordId) { word in
                        WordCardRow(word: word)
                    }
                }
                .padding()
            }
            .scrollIndicators(.hidden)
            .scrollDismissesKeyboard(.immediately)
        }
        .onChange(of: prefix) {
            presenter.filter(prefix: prefix)
        }
        .modifier(SheetViewModifier())
    }
}

struct DictionaryPreview: PreviewProvider {

    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        var body: some View {
            DictionaryView(
                isLearned: true
            )
        }
    }
}

