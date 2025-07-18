//
//  AIChatView.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation
import SwiftUI

struct AIChatView: View {
    // MARK: - Properties
    @State var inputMessage: String = ""
    @State var showExitAlert = false
    @FocusState var textFieldFocused

    var endId = UUID()

    var onExit: (() -> Void)?

    @ObservedObject var presenter: AIChatScreenPresenter

    // MARK: - View Constances
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)
        static let backButtonPadding = CGFloat(35)
        // Corner Radiuses
        static let sheetCornerRadius = CGFloat(20)
        static let gridInfoVerticalPadding = CGFloat(20)
    }

    var body: some View {
        ZStack {
            messagesGrid
            inputView
            backButton
                .frame(
                    maxWidth: .infinity,
                    maxHeight: .infinity,
                    alignment: .topLeading
                )
                .padding(Const.backButtonPadding)
                .ignoresSafeArea()
                .onTapGesture {
                    showExitAlert = true
                }
        }
        .onAppear {
            presenter.sendMessage("Hello!")
        }
        .alert("Are you sure, that you want exit?", isPresented: $showExitAlert) {
            Button ("No", role: .cancel) {
                showExitAlert = false
            }
            Button ("Yes", role: .destructive) {
                presenter.finishChat()
                onExit?()
            }
        }
    }

    // MARK: - Subviews

    private var messagesGrid: some View {
        ScrollViewReader { proxy in
            ScrollView {
                VStack {
                    ForEach(presenter.messages, id: \.text) { message in
                        MessageView(
                            text: message.text,
                            role: message.role
                        )
                    }
                    Spacer(
                        minLength: 60
                    )
                    Text("")
                        .id(endId)
                }
                .scrollIndicators(.hidden)
                .onAppear {
                    withAnimation {
                        proxy.scrollTo(endId, anchor: .bottom)
                    }
                }
                .onChange(of: presenter.messages.count) { _, _ in
                    withAnimation {
                        proxy.scrollTo(endId, anchor: .bottom)
                    }
                }
            }
            .scrollDismissesKeyboard(.interactively)
        }
    }

    /// Input field + button to send message
    private var inputView: some View {
        HStack(
            alignment: .bottom,
            spacing: 12
        ){
            inputField
            sendButton {
                sendMessage()
            }
        }
        .scrollDismissesKeyboard(.interactively)
        .padding(.horizontal, 12)
        .padding(.top, 8)
        .safeAreaPadding(8)
        .background(
            VStack {
                Spacer()
                RoundedRectangle(cornerRadius: 20)
                    .fill(.clear)
                    .glass(
                        cornerRadius: 20,
                        opacity: 0.8
                    )
                    .ignoresSafeArea()
                    .frame(
                        alignment: .bottom
                    )
            }
        )
        .frame(
            maxHeight: .infinity,
            alignment: .bottom
        )
    }

    /// input field
    private var inputField: some View {
        TextField("Message", text: $inputMessage, axis: .vertical)
            .lineLimit(5, reservesSpace: false)
            .frame(maxWidth: .infinity)
            .padding()
            .clipShape(
                RoundedRectangle(cornerRadius: 20)
            )
            .background(
                RoundedRectangle(cornerRadius: 20)
                    .fill(.orangeSecondary)
                    .stroke(.orangePrimary, lineWidth: 2)
            )
            .focused($textFieldFocused)
            .onSubmit(of: .text) {
                inputMessage.append("\n")
                textFieldFocused = true
            }
    }

    private var backButton: some View {
        Image(systemName: "arrowshape.turn.up.left.fill")
            .font(.title3)
            .foregroundStyle(.whiteBackground)
            .padding()
            .glass(
                cornerRadius: 100,
                fill: .orangePrimary,
                opacity: 0.8
            )
    }

    /// button to send message
    private func sendButton(action: @escaping () -> Void) -> some View {
        Button {
            action()
        } label: {
            Image(systemName: "arrowshape.up.circle.fill")
                .font(.largeTitle)
                .foregroundStyle(.orangePrimary)
                .padding(.bottom, 10)
        }
    }

    /// model of message
    private func mesageView(message: MessageModel) -> some View {
        Text(message.text)
            .frame(alignment: message.role == .ai ? .leading : .trailing)
    }

    /// send message in presenter
    private func sendMessage() {
        guard !inputMessage.isEmpty else { return }

        presenter.sendMessage(inputMessage)
        inputMessage = ""
    }
}


// MARK: - Preview
struct AIChatPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        var body: some View {
            AIChatBuilder.build() {
                print("exit")
            }
        }
    }
}
