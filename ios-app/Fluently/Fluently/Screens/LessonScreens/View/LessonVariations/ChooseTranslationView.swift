//
//  ChooseTranslationView.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftUI

struct ChooseTranslationView: View {
    // MARK: - Properties
    @State var word: String
    @State var selectedAnswer: String?
    @State var answers: [String]

    @State var correctAnswer: String

    @State var isCorrect = false
    @State var answerIsShown = false

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

    // MARK: - Subviews
    var buttonNext: some View {
        Button {
            if let selectedAnswer {
                if answerIsShown {
                    onAnswerSelected(selectedAnswer)
                } else {
                    withAnimation(.easeIn(duration: 0.3)) {
                        answerIsShown = true
                        isCorrect = correctAnswer == selectedAnswer
                    }
                }
//                onAnswerSelected(selectedAnswer)
            }
        } label: {
            Text("Next")
                .padding()
                .frame(maxWidth: .infinity)
                .massiveButton(color: .blue)
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
                    isCorrectAnswer: answer == correctAnswer,
                    isSubmitted: answerIsShown,
                    answer: answer
                ) {
                        selectedAnswer = answer
                    }
            }
        }
    }
}

