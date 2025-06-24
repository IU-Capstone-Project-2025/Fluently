//
//  PickOptionsView.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI


struct PickOptionsView: View {
    @State var sentence: String
    @State var selectedAnswer: String?
    @State var answers: [String]

    var onAnswerSelected: (String) -> Void

    var body: some View {
        VStack {
            sentenve()
                .foregroundStyle(.blackText)
                .padding()
            listOfAnswers

            Spacer()
            
            buttonNext
        }
    }

    var buttonNext: some View {
        Button {
            if let selectedAnswer {
                onAnswerSelected(selectedAnswer)
            }
        } label: {
            ZStack(alignment: .center) {
                RoundedRectangle(cornerRadius: 12)
                    .fill(.blueSecondary)
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(.blueAccent, lineWidth: 2)
                    )
                RoundedRectangle(cornerRadius: 12)
                    .fill(.blueAccent)
                    .frame(height: 6)
                    .offset(y: 6)
                    .mask(RoundedRectangle(cornerRadius: 12))

                Text("Next")
                    .foregroundStyle(.blueAccent)
                    .padding(.vertical, 12)
            }
            .grayscale( selectedAnswer == nil ? 1 : 0)
//            .frame(maxWidth: .infinity)
            .frame(maxHeight: 60)
            .padding(.horizontal, 100)
        }
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

    func sentenve() -> some View {
        var parts = sentence.split(separator: "____")

        if let selectedAnswer {
            return Text("\(parts[0]) _\(selectedAnswer)_ \(parts[1])")
        }
        return Text("\(parts[0]) ____ \(parts[1])")
    }
}

struct AnswerButton: View {
    var isSelected: Bool
    var answer: String

    var onTap: () -> Void

    var body: some View {
        Button {
            onTap()
        } label: {
            ZStack(alignment: .center) {
                RoundedRectangle(cornerRadius: 12)
                    .fill(.orangeSecondary)
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(.orangePrimary, lineWidth: 2)
                    )
                RoundedRectangle(cornerRadius: 12)
                    .fill(.orangePrimary)
                    .frame(height: 6)
                    .offset(y: 6)
                    .mask(RoundedRectangle(cornerRadius: 12))

                Text(answer)
                    .foregroundStyle(.orangePrimary)
                    .padding(.vertical, 12)
            }
            .grayscale( isSelected ? 1 : 0)
//            .frame(maxWidth: .infinity)
            .frame(maxHeight: 60)
            .padding(.horizontal, 100)
        }
    }
}
