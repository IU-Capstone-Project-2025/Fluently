//
//  LessonInfo.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI

struct LessonInfo: View {
    // MARK: - Properties
    let minutes: Int
    let seconds: Int

    var body: some View {
        HStack{
            learnImage
            VStack(alignment: .leading) {
                Text("Let's learn how to speak fluently!")
                    .foregroundStyle(.blackText)
                    .font(.appFont.callout)
                Text("\(minutes) minutes \(seconds) seconds")
                    .foregroundStyle(.blackText.opacity(0.6))
                    .font(.appFont.subheadline)
            }
        }
        .padding(4)
        .background(
            RoundedRectangle(cornerRadius: 50)
                .fill(.grayFluently.opacity(0.4))
        )
    }

    // MARK: - Subviews
    
    var learnImage: some View {
        Image(systemName: "graduationcap.fill")
            .foregroundStyle(.blueAccent)
            .padding()
            .background(
                Circle()
                    .fill(.whiteBackground)
            )
    }
}
