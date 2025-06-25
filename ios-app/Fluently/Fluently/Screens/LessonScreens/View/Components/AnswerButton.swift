//
//  AnswerButton.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//


import SwiftUI

struct AnswerButton: View {
    var isSelected: Bool
    var answer: String

    var onTap: () -> Void

    var body: some View {
        Button {
            onTap()
        } label: {
            Text(answer)
                .foregroundStyle(.orangePrimary)
                .padding()
                .frame(maxWidth: .infinity)
                .modifier(ButtonViewModifier(color: .orange))
                .grayscale( isSelected ? 1 : 0)
                .frame(maxHeight: 60)
        }
        .buttonStyle(PlainButtonStyle())
    }
}
