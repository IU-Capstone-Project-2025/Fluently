//
//  WordCardView.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftUI

// TODO: specifi the min width and min height 
struct WordCardView: View {
    // MARK: - Properties
    @State var word: WordModel

    var onKnowTapped: () -> Void
    var onLearnTapped: () -> Void

    var body: some View {
        VStack {

            wordCard

            Spacer()

            HStack (spacing: 30) {
                buttonKnow
                buttonLearn
            }
            .padding(.horizontal, 30)
            .padding(.bottom, 100)

        }
    }

    // MARK: - Subviews

    // card word representation
    var wordCard: some View{
        VStack(alignment: .leading, spacing: 10) {
            VStack (alignment: .leading, spacing: 4) {
                Text(word.word)
                    .foregroundStyle(.blackText)
                    .font(.appFont.title)
                Text(word.translation)
                    .foregroundStyle(.blackText.opacity(0.6))
                    .font(.appFont.subheadline)
                Text(word.transcription)
                    .foregroundStyle(.blackText)
                    .font(.appFont.subheadline)
            }

            Text("\(word.topic) : \(word.subtopic)")
                .foregroundStyle(.blackText)
                .font(.appFont.caption)

            VStack(alignment: .leading, spacing: 4) {
                ForEach(word.sentences, id: \.self) { sentence in
                    VStack(alignment: .leading) {
                        Text("- \(sentence.text)")
                            .foregroundColor(.blackText)
                        Text("- \(sentence.translation)")
                            .font(.appFont.caption)
                            .foregroundColor(.blackText.opacity(0.6))
                    }
                }
            }
        }
        .padding()
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(.blueSecondary)
                .stroke(.blueAccent, lineWidth: 2)
        )
    }

    var buttonLearn: some View {
        Button {
            onLearnTapped()
        } label: {
            Text("Learn")
                .padding()
                .frame(maxWidth: .infinity)
                .massiveButton(color: .blue)
                .frame(maxHeight: 60)
        }
        .buttonStyle(PlainButtonStyle())
    }

    var buttonKnow: some View {
        Button {
            onKnowTapped()
        } label: {
            Text("Know")
                .padding()
                .frame(maxWidth: .infinity)
                .massiveButton(color: .blue)
                .frame(maxHeight: 60)
        }
        .buttonStyle(PlainButtonStyle())
    }
}
