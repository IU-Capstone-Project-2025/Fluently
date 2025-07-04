//
//  TypeTranslationView.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftUI

struct TypeTranslationView: View {
    // MARK: - Properties
    @State var word: String
    @State var typedAnswer: String = ""

    var onAnswerSelected: (String) -> Void

    var body: some View {
        ScrollView{
            Text(word)
                .foregroundStyle(.blackText)
                .font(.appFont.title)
                .padding()

            answerField
                .padding(.vertical)
                .padding(.horizontal, 100)

            Spacer()

            buttonNext
                .padding(.horizontal, 100)

            Spacer()
        }
        .scrollIndicators(.hidden)
        .scrollDismissesKeyboard(.interactively)
        .ignoresSafeArea(.keyboard)
    }

    // MARK: - Subviews

    var buttonNext: some View {
        Button {
            if !typedAnswer.isEmpty {
                onAnswerSelected(typedAnswer.lowercased().trimmingCharacters(in: .whitespacesAndNewlines))
            }
        } label: {
            Text("Next")
                .padding()
                .frame(maxWidth: .infinity)
                .massiveButton(color: .blue)
                .grayscale( typedAnswer.isEmpty ? 1 : 0)
                .frame(maxHeight: 60)
        }
        .buttonStyle(PlainButtonStyle())
    }

    var answerField: some View {
        VStack (alignment: .center, spacing: 10) {
            TextField("Type translation", text: $typedAnswer)
                .lineLimit(1)
                .padding()
                .frame(maxWidth: .infinity)
                .massiveButton(color: .orange)
                .frame(maxHeight: 60)
        }
    }
}
