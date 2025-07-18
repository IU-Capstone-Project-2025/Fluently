//
//  ButtonViewModifier.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import SwiftUI

// MARK: - View modifier
// Modifiier for 3D Button
struct ButtonViewModifier: ViewModifier {
    // MARK: - Color
    enum ButtonColor {
        case blue
        case orange
        case purple
        case red
        case green
    }

    // MARK: - Properties
    var color: ButtonColor = .orange
    private var primaryColor: Color
    private var secondaryColor: Color

    @State var isPressed = false

    // MARK: - Init
    init(color: ButtonColor) {
        self.color = color

        switch color {
            case .orange:
                primaryColor = .orangePrimary
                secondaryColor = .orangeSecondary
            case .blue:
                primaryColor = .blueAccent
                secondaryColor = .blueSecondary
            case .purple:
                primaryColor = .purpleAccent
                secondaryColor = .purpleSecondary
            case .red:
                primaryColor = .redPrimary
                secondaryColor = .redSecondary
            case .green:
                primaryColor = .greenPrimary
                secondaryColor = .greenSecondary
        }
    }

    func body(content: Content) -> some View {
        content
            .foregroundStyle(primaryColor)
            .offset(y: isPressed ? 6 : 2)
            .background(
                ZStack(alignment: .center) {
                    VStack {
                        Spacer()
                        RoundedRectangle(cornerRadius: 12)
                            .fill(primaryColor)
                            .stroke(primaryColor, lineWidth: 1)
                            .frame(alignment: .bottom)
                            .offset(y: isPressed ? 10 : 6)
                    }

                    RoundedRectangle(cornerRadius: 12)
                        .fill(secondaryColor)
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(primaryColor, lineWidth: 2)
                        )
                        .offset(y: isPressed ? 4 : 0)
                }
            )
            .animation(.easeOut(duration: 0.1), value: isPressed)
            .simultaneousGesture(
                DragGesture(minimumDistance: 0)
                    .onChanged { _ in isPressed = true }
                    .onEnded { _ in isPressed = false }
            )
    }
}

// MARK: - View extension
extension View {
    func massiveButton(color: ButtonViewModifier.ButtonColor) -> some View {
        modifier(ButtonViewModifier(color: color))
    }
}

// MARK: - Preview
struct ButtonsPreview: PreviewProvider {

    static var previews: some View {
        VStack(spacing: 20) {
            Text("Orange")
                .padding()
                .massiveButton(color: .orange)
            Text("Blue")
                .padding()
                .massiveButton(color: .blue)
            Text("Purple")
                .padding()
                .massiveButton(color: .purple)
            Text("Red")
                .padding()
                .massiveButton(color: .red)
            Text("Green")
                .padding()
                .massiveButton(color: .green)
        }
    }
}
