//
//  AnswerButton.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//


import SwiftUI

struct AnswerButton: View {
    // MARK: - Properties
    var isSelected: Bool
    var isCorrectAnswer: Bool
    var isSubmitted: Bool
    var answer: String

    var onTap: () -> Void

    var body: some View {
        Text(answer)
            .padding()
            .frame(maxWidth: .infinity)
            .massiveButton(color: isSubmitted ? isCorrectAnswer ? .green : .red : .orange)
            .grayscale( isSubmitted && (isSelected || isCorrectAnswer) ? 0 : isSelected ? 0 : 1)
            .frame(maxHeight: 60)
            .onTapGesture {
                onTap()
            }
            .disabled(isSubmitted)
    }
}

struct OpaqueButtonStyle: ButtonStyle {
    func makeBody(configuration: Configuration) -> some View {
        configuration.label
            .opacity(1.0)
    }
}
