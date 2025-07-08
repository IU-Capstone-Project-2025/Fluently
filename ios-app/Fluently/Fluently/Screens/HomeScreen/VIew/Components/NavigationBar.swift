//
//  NavigationBar.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import SwiftUI

struct NavigationBar: View {

    enum Screen: String, CaseIterable {
        case calendar = "calendar"
        case home = "house.fill"
        case statistic = "chart.bar"

        var title: String {
            switch self {
                case .calendar: return "Calendar"
                case .home: return "Home"
                case .statistic: return "Stats"
            }
        }
    }

    @Binding var currentScreen: Screen

    var body: some View {
        HStack {
            ForEach(Screen.allCases, id: \.self) { screen in
                MenuButton(
                    isSelected: currentScreen == screen,
                    imageName: screen.rawValue,
                    name: screen.title,
                    onSelect: {
                        currentScreen = screen
                    }
                )
            }
            .padding(.horizontal)
        }
        .padding()
    }
}

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


struct NavigationPreview: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        @State var screen = NavigationBar.Screen.home
        var body: some View {
            NavigationBar(currentScreen: $screen)
        }
    }
}
