//
//  Smegma.swift
//  Fluently
//
//  Created by Савва Пономарев on 30.06.2025.
//

import SwiftUI

struct Smegma: Shape {
    var percent: Double

    func path(in rect: CGRect) -> Path {
        Path { path in
            path.move(to: CGPoint(x: rect.minX, y: rect.minY))                              /// left-top cornet
            path.addLine(to: CGPoint(x: rect.midX, y: rect.minY + percent * rect.height))   /// top center
            path.addLine(to: CGPoint(x: rect.maxX, y: rect.minY))                           /// right-top corner
            path.addLine(to: CGPoint(x: rect.maxX - percent * rect.width, y: rect.midY))    /// right cener
            path.addLine(to: CGPoint(x: rect.maxX, y: rect.maxY))                           /// right-botton corner
            path.addLine(to: CGPoint(x: rect.midX, y: rect.maxY - percent * rect.height))   /// bottom center
            path.addLine(to: CGPoint(x: rect.minX, y: rect.maxY))                           /// left-bottom corner
            path.addLine(to: CGPoint(x: rect.minX + percent * rect.width, y: rect.midY))    /// left center
            path.closeSubpath()
        }
    }
}

struct SmegmaRounded: Shape {
    var percent: Double

    func path(in rect: CGRect) -> Path {
        Path { path in
            path.move(to: CGPoint(x: rect.midX, y: rect.minY + percent * rect.height))

            path.addCurve(
                to: CGPoint(x: rect.maxX - percent * rect.width, y: rect.midY),
                control1:  CGPoint(x: rect.maxX, y: rect.minY),
                control2:  CGPoint(x: rect.maxX + percent * rect.width, y: rect.midY)
            ) /// Top Petal

            path.addCurve(
                to: CGPoint(x: rect.midX, y: rect.maxY - percent * rect.height),
                control1:  CGPoint(x: rect.maxX, y: rect.maxY),
                control2:  CGPoint(x: rect.midX, y: rect.maxY + percent * rect.height)
            ) /// Right Petal

            path.addCurve(
                to: CGPoint(x: rect.minX + percent * rect.width, y: rect.midY),
                control1: CGPoint(x: rect.minX, y: rect.maxY),
                control2: CGPoint(x: rect.minX - percent * rect.width, y: rect.midY)
            ) /// Bottom Petal

            path.addCurve(
                to: CGPoint(x: rect.midX, y: rect.minY + percent * rect.height),
                control1: CGPoint(x: rect.minX, y: rect.minY),
                control2: CGPoint(x: rect.midX, y: rect.minY - percent * rect.height)
            ) /// Left Petal

            path.closeSubpath()
        }
    }
}
