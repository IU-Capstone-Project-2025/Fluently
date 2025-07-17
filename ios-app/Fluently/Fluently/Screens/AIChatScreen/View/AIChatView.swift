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
    @FocusState var textFieldFocused

    var endId = UUID()

    @State var messages: [MessageModel] = MessageModel.mockGenerator()

    // MARK: - View Constances
    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)

        // Corner Radiuses
        static let sheetCornerRadius = CGFloat(20)
        static let gridInfoVerticalPadding = CGFloat(20)
    }

    var body: some View {
        ZStack {
            messagesGrid
            inputView
        }
    }

    // MARK: - Subviews

    private var messagesGrid: some View {
        ScrollViewReader { proxy in
            ScrollView {
                VStack {
                    ForEach(messages, id: \.text) { message in
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
                .onChange(of: messages.count) { _, _ in
                    withAnimation {
                        proxy.scrollTo(endId, anchor: .bottom)
                    }
                }
                .scrollDismissesKeyboard(.immediately)
            }
        }
    }

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
        .padding(.horizontal, 12)
        .padding(.top, 20)
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

    private var inputField: some View {
        TextField("Messaege", text: $inputMessage, axis: .vertical)
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

    private func mesageView(message: MessageModel) -> some View {
        Text(message.text)
            .frame(alignment: message.role == .ai ? .leading : .trailing)
    }

    private func sendMessage() {
        guard !inputMessage.isEmpty else { return }

        let newMessage = MessageModel(text: inputMessage, role: .user)
        messages.append(newMessage)
        inputMessage = ""

        DispatchQueue.main.asyncAfter(deadline: .now() + 1) {
            let aiResponse = MessageModel(text: "This is an automated response", role: .ai)
            messages.append(aiResponse)
        }
    }
}


// MARK: - Preview
struct AIChatPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        var body: some View {
            AIChatView()
        }
    }
}
