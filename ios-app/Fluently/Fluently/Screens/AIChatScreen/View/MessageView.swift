//
//  MessageView.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import SwiftUI

struct MessageView: View {
    var text: String
    var role: MessageRole

    var body: some View {
        HStack {
            if role == .ai {
                content
                Spacer(minLength: 60)
            } else {
                Spacer(minLength: 60)
                content
            }
        }
        .padding(.horizontal, 8)
    }

    private var content: some View {
        Text(LocalizedStringKey(stringLiteral: text))
            .foregroundColor(.blackText)
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(role == .ai ? .orangePrimary : .orangeSecondary)
            )
            .frame(maxWidth: .infinity, alignment: role == .ai ? .leading : .trailing)
            .fixedSize(horizontal: false, vertical: true)
    }
}
