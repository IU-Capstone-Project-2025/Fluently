//
//  GlassView.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import SwiftUI

// view to make items look like glass
struct GlassView: View {
    let cornerRadius: CGFloat
    let fill: Color
    let opacity: CGFloat
    let shadowRadius: CGFloat

    init(cornerRadius: CGFloat, fill: Color = .white, opacity: CGFloat = 0.25, shadowRadius: CGFloat = 10.0) {
        self.cornerRadius = cornerRadius
        self.fill = fill
        self.opacity = opacity
        self.shadowRadius = shadowRadius
    }

    var body: some View {
        RoundedRectangle(cornerRadius: cornerRadius)
            .fill(fill)
            .opacity(opacity)
            .shadow(radius: shadowRadius)
    }
}

// MARK: - View modifier
struct GlassModifier: ViewModifier {
    let cornerRadius: CGFloat
    let fill: Color
    let opacity: CGFloat
    let shadowRadius: CGFloat

    func body(content: Content) -> some View {
        content
            .background {
                GlassView(cornerRadius: cornerRadius, fill: fill, opacity: opacity, shadowRadius: shadowRadius)
            }
    }
}

// MARK: - View extension
extension View {
    func glass(
        cornerRadius: CGFloat,
        fill: Color = .white,
        opacity: CGFloat = 0.25,
        shadowRadius: CGFloat = 10.0
    ) -> some View {
        modifier(GlassModifier(cornerRadius: cornerRadius, fill: fill, opacity: opacity, shadowRadius: shadowRadius))
    }
}
