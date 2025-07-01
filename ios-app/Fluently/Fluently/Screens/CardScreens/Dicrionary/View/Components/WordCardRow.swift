//
//  WordCardRow.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import SwiftUI

struct WordCardRow: View {
    let word: Word
    @State var isHidden = true

    var body: some View {
        HStack(alignment: .center) {
            VStack {
                Spacer()
                    .fixedSize()
                audio
                    .fixedSize()
                Spacer()
            }

            VStack(alignment: .leading, spacing: 4) {
                textInfo
                if !isHidden {
                    VStack (alignment: .leading, spacing: 4){
                        transcription
                        sentences
                            .padding(.vertical, 3)
                    }
                    .transition(.move(edge: .top).combined(with: .opacity))
                }
            }
        }
        .padding(.horizontal)
        .frame(maxWidth: .infinity, alignment: .leading)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(.orangeSecondary)
                .stroke(.orangePrimary, lineWidth: 2)
        )
        .onTapGesture {
            withAnimation (.easeInOut(duration: 0.3)) {
                isHidden.toggle()
            }
        }
    }

    var audio: some View {
        Image(systemName: "speaker.wave.2")
            .font(.title)
            .foregroundStyle(.blueAccent)
    }

    var textInfo: some View {
        VStack(alignment: .leading) {
            Text(word.word)
                .font(.appFont.title2)
                .foregroundStyle(.blackText)
            Text(word.translation)
                .font(.appFont.subheadline)
                .foregroundStyle(.grayFluently)

        }
    }

    var transcription: some View {
        Text(word.transcription)
            .font(.appFont.title3)
            .foregroundStyle(.blackText)
    }

    var sentences: some View {
        VStack(alignment: .leading, spacing: 6) {
            ForEach(word.sentences, id: \.self) { sentence in
                Text("- \(sentence)")
                    .font(.appFont.headline)
                    .foregroundStyle(.blackText)
            }
        }
    }
}
