//
//  PickOptionsView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct PickOptionsView: View {
    // MARK: - Properties
    @State var sentence: String
    @State var selectedAnswer: String?
    @State var answers: [String]

    @State var correctAnswer: String

    @State var isCorrect = false
    @State var answerIsShown = false

    var onAnswerSelected: (String) -> Void

    var body: some View {
        VStack {
            sentenceView()
                .foregroundStyle(.blackText)
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
        Text("Next")
            .padding()
            .frame(maxWidth: .infinity)
            .massiveButton(color: .blue)
            .grayscale( selectedAnswer == nil ? 1 : 0)
            .frame(maxHeight: 60)
            .onTapGesture {
                if let selectedAnswer {
                    if answerIsShown {
                        onAnswerSelected(selectedAnswer)
                    } else {
                        withAnimation(.easeIn(duration: 0.3)) {
                            answerIsShown = true
                            isCorrect = correctAnswer == selectedAnswer
                        }
                    }
                }
            }
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

    // updating sentence view
    func sentenceView() -> some View {
        let parts = sentence.split(separator: "_")

        if let selectedAnswer {
            if parts.count == 2 {
                return Text("\(parts[0]) _\(selectedAnswer)_ \(parts[1])")
            } else {
                return Text("\(parts[0]) _\(selectedAnswer)_")
            }
        }
        if parts.count == 2 {
            return Text("\(parts[0]) __________ \(parts[1])")
        } else {
            return Text("\(parts[0]) __________")
        }
    }
}
