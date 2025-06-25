//
//  ChooseTranslationView.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftUI

struct ChooseTranslationView: View {
    @State var word: String
    @State var selectedAnswer: String?
    @State var answers: [String]

    var onAnswerSelected: (String) -> Void

    var body: some View {
        VStack {
            Text(word)
                .foregroundStyle(.blackText)
                .font(.appFont.title)
                .padding()
            listOfAnswers
                .padding(.horizontal, 100)

            Spacer()

            buttonNext
                .padding(.horizontal, 100)

            Spacer()
        }
    }

    var buttonNext: some View {
        Button {
            if let selectedAnswer {
                onAnswerSelected(selectedAnswer)
            }
        } label: {
            Text("Next")
                .padding()
                .frame(maxWidth: .infinity)
                .modifier(ButtonViewModifier(color: .blue))
                .grayscale( selectedAnswer == nil ? 1 : 0)
                .frame(maxHeight: 60)
        }
        .buttonStyle(PlainButtonStyle())
    }

    var listOfAnswers: some View {
        VStack (alignment: .center, spacing: 10) {
            ForEach(answers, id: \.self ) { answer in
                AnswerButton (
                    isSelected: selectedAnswer == answer,
                    answer: answer) {
                        selectedAnswer = answer
                    }
            }
        }
    }
}

