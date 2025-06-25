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
        Button {
            if let selectedAnswer {
                onAnswerSelected(selectedAnswer)
            }
        } label: {
            Text("Next")
                .padding()
                .frame(maxWidth: .infinity)
                .foregroundStyle(.blueAccent)
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

    // updating sentence view
    func sentenceView() -> some View {
        let parts = sentence.split(separator: "____")

        if let selectedAnswer {
            return Text("\(parts[0]) _\(selectedAnswer)_ \(parts[1])")
        }
        return Text("\(parts[0]) ____ \(parts[1])")
    }
}
