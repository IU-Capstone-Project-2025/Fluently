//
//  WordCardView.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation
import SwiftUI

struct WordCardView: View {
    // MARK: - Properties
    @State var word: WordModel

    var onKnowTapped: () -> Void
    var onLearnTapped: () -> Void

    var body: some View {
        VStack {

            wordCard
                .padding(20)

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
                Text(word.word!)
                    .foregroundStyle(.blackText)
                    .font(.appFont.title)
                Text(word.translation!)
                    .foregroundStyle(.blackText.opacity(0.6))
                    .font(.appFont.subheadline)
                Text(word.transcription!)
                    .foregroundStyle(.blackText)
                    .font(.appFont.subheadline)
            }

            Text("\(word.topic!) : \(word.subtopic!)")
                .foregroundStyle(.blackText)
                .font(.appFont.caption)

            VStack(alignment: .leading, spacing: 4) {
                ForEach(word.sentences!.prefix(3), id: \.self) { sentence in
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
        .frame(
            minWidth: 120,
            minHeight: 120
        )
        .padding()
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(.blueSecondary)
                .stroke(.blueAccent, lineWidth: 2)
        )
    }

    var buttonLearn: some View {
        Text("Learn")
            .padding()
            .frame(maxWidth: .infinity)
            .massiveButton(color: .blue)
            .frame(maxHeight: 60)
            .onTapGesture {
                onLearnTapped()
            }
    }

    var buttonKnow: some View {
        Text("Know")
            .padding()
            .frame(maxWidth: .infinity)
            .massiveButton(color: .blue)
            .frame(maxHeight: 60)
            .onTapGesture {
                onKnowTapped()
            }
    }
}

// MARK: Preview Provider
struct WordModelPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        let word = WordModel.mockWord()
        var body: some View {
            WordCardView(
                word: word,
                onKnowTapped: {},
                onLearnTapped: {}
            )
        }
    }
}
