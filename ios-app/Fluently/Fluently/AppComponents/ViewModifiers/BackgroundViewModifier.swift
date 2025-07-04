//
//  BackgroundViewModifier.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

// Full orange view modifier
struct BackgroundViewModifier: ViewModifier {

    var colorful: Bool = false

    func body(content: Content) -> some View {
        content
            .containerRelativeFrame([.horizontal, .vertical])
            .background(backgroundContent)
            .background(.orangePrimary)
    }

    @ViewBuilder
    private var backgroundContent: some View {
        if colorful {
            ZStack {
                colorfulBackground
            }
        } else {
            Color.orangePrimary
        }
    }

    private var colorfulBackground: some View {
        GeometryReader { proxy in
            let size = proxy.size
            Circle()
                .fill(.pink)
                .padding(20)
                .blur(radius: 120)
                .offset(x: size.width / 1.8)
            Circle()
                .fill(.purpleSecondary)
                .padding(20)
                .blur(radius: 100)
                .offset(
                    x: size.width / 3,
                    y: size.height / 2
                )
            Circle()
                .fill(.orangeSecondary)
                .padding(30)
                .blur(radius: 120)
                .offset(
                    x: -size.width * 0.15,
                    y: size.height / 1.2
                )
            Circle()
                .fill(.blueSecondary)
                .blur(radius: 120)
                .offset(
                    x: -size.width * 0.7,
                    y: size.height * 0.2
                )
        }
    }
}

struct BackgroundPreview: PreviewProvider {
    static var previews: some View {
        Text("preview")
            .modifier(BackgroundViewModifier(colorful: true))
    }
}
