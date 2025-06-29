//
//  LaunchScreenView.swift
//  Fluently
//
//  Created by Савва Пономарев on 29.06.2025.
//

import Foundation
import SwiftUI

struct LaunchScreenView: View {
    @Binding var isActive: Bool
    
    @State private var size = 0.8
    @State private var opacity = 0.5

    @State private var rotateAngle1: Double = 40
    @State private var rotateAngle2: Double = 134

    var body: some View {
        ZStack {
//            Smegma(percent: 0.2)
//                .fill(.blueAccent.opacity(0.4))
//                .frame(width: 200 * size, height: 160 * size)
//                .rotationEffect(Angle(degrees: rotateAngle2))
//                .blur(radius: 25)
//            Smegma(percent: 0.1)
//                .fill(.orange.opacity(0.6))
//                .frame(width: 200 * size, height: 200 * size)
//                .rotationEffect(Angle(degrees: rotateAngle1))
//                .blur(radius: 25)
            Text("FLUENTLY")
                .font(.appFont.largeTitle)
                .foregroundStyle(.orangePrimary)
                .scaleEffect(size)
                .opacity(opacity)
                .onAppear{
                    withAnimation(.easeIn(duration: 1.2)) {
                        self.opacity = 0.9
                        self.size = 1.5
                        self.rotateAngle1 = 140
                        self.rotateAngle2 = 201
                    }
                }
        }
        .containerRelativeFrame([.horizontal, .vertical])
        .background(.whiteBackground)
        .onAppear {
            DispatchQueue.main.asyncAfter(deadline: .now() + 1.5){
                self.isActive = false
            }
        }
    }
}

struct LaunchScreen: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        @State var isActive = true
        var body: some View {
            LaunchScreenView(isActive: $isActive)
        }
    }
}
