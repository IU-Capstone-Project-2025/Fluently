//
//  NotesView.swift
//  Fluently
//
//  Created by Савва Пономарев on 30.06.2025.
//

import Foundation
import SwiftUI

struct NotesView: View {
    @ObservedObject var presenter: NotesScreenPresenter

    @Environment(\.dismiss) var dismiss

    // MARK: - Constants
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)

        // Sizes
        static let recordButtonSize = CGFloat(75)
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
            Text("Notes")
                .foregroundStyle(.whiteText)
                .font(.appFont.largeTitle.bold())
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(.horizontal, Const.horizontalPadding)
        }
    }

    var recordingButton: some View {
        Button {
            withAnimation() {
                presenter.toggleRecording()
            }
        } label: {
            ZStack {
                RoundedRectangle(cornerRadius: presenter.isRecording ? 12 : 40)
                    .fill(.orangePrimary)
                    .padding(presenter.isRecording ? 15 : 4)
                Circle()
                    .fill(.clear)
                    .stroke(
                        .grayFluently,
                        lineWidth: 4
                    )
            }
        }
    }

    ///  Grid with main info
    var infoGrid: some View {
        VStack (alignment: .center) {

            Spacer()

            recordingButton
                .frame(
                    maxWidth: Const.recordButtonSize,
                    maxHeight: Const.recordButtonSize,
                    alignment: .bottom
                )
                .padding(.bottom, 30)
        }
        .modifier(SheetViewModifier())
    }
}


struct NotesPreview: PreviewProvider {

    static var previews: some View {
        NotesPreviewWrapper()
    }

    struct NotesPreviewWrapper: View {
        let router: NotesScreenRouter
        let presenter: NotesScreenPresenter

        init() {
            self.router = NotesScreenRouter(
                router: AppRouter()
            )
            self.presenter = NotesScreenPresenter(
                router: router
            )
        }

        var body: some View {
            NotesView(
                presenter: presenter
            )
        }
    }
}
