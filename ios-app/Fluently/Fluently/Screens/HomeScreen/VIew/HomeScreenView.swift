//
//  HomeScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

struct HomeScreenView: View {

    var body: some View {
        VStack {
            welcomText
            ZStack {
                statistic
                    .frame(height: 140)
                HStack {
                    Spacer()
                    startButton
                        .frame(width: 120, alignment: .trailing)
                        .padding(.trailing, 20)
                }
            }
        }
        .padding()
        .containerRelativeFrame([.horizontal, .vertical])
        .background(Color.theme.backgroundColor)
    }

    private var welcomText: some View {
        VStack (alignment: .leading) {
            Text("Good afternoon User!")
                .font(.headline)
                .foregroundStyle(Color.theme.primary)
                .padding(.horizontal)
            Text("Lorem Ipsum is simply dummy text of the printing and typesetting industry.")
                .font(.title)
                .foregroundStyle(.white)
                .padding(.horizontal)

        }
    }

    private var statistic: some View {
        RoundedRectangle(cornerRadius: 80)
            .fill(Color.theme.complementary2)
    }

    private var startButton: some View {
        ZStack () {
            Circle()
                .fill(Color.theme.complementary1)
            Image(systemName: "arrow.right")
                .font(.largeTitle)
        }
    }
}

#Preview {
    HomeScreenView()
}
