//
//  MenuButton.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import SwiftUI

struct MenuButton: View {

    var isSelected: Bool = false

    var imageName: String = "house.fill"
    var name: String = "Home"

    var onSelect: () -> Void

    var body: some View {
        HStack {
            Image(systemName: imageName)
                .font(.headline)
                .foregroundStyle(.orangePrimary)
            nameLabel
        }
        .onTapGesture {
            withAnimation(
                .interpolatingSpring(
                    mass: 1.0,
                    stiffness: 150,
                    damping: 16.5,
                    initialVelocity: 0
                )
            ) {
                onSelect()
            }
        }
        .padding(.horizontal, 8)
        .padding(.vertical, 6)
        .glass(cornerRadius: 50, fill: .orangePrimary)
        .scaleEffect(isSelected ? 1.05 : 1.0)
    }

    @ViewBuilder
    var nameLabel: some View {
        if isSelected {
            Text(name)
                .foregroundStyle(.orangePrimary)
                .font(.appFont.secondaryHeadline)
        }
    }
}
