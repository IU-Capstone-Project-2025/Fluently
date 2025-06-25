//
//  ScreenCard.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI

struct ScreenCard: View {
    // MARK: - Properties
    private let type: CardType
    private let onTap: () -> Void
    private let count: String

    // MARK: - Initializer
    init(type: CardType, count: String = "0", onTap: @escaping () -> Void) {
        self.type = type
        self.count = count
        self.onTap = onTap
    }

    // MARK: - Body
    var body: some View {
        content
            .background(background)
            .contentShape(Rectangle())
            .onTapGesture(perform: onTap)
    }

    // MARK: - Subviews
    private var content: some View {
        VStack(alignment: .leading, spacing: 4) {
            Image(systemName: type.imageName)
                .foregroundStyle(type.color.primary)

            Text(count)
                .foregroundStyle(.blackText)
                .font(.appFont.title)
                .padding(.top, 6)

            Text(type.rawValue)
                .foregroundStyle(.blackText)
                .font(.appFont.caption2)
        }
        .frame(maxWidth: .infinity, alignment: .leading)
        .padding()
    }

    private var background: some View {
        RoundedRectangle(cornerRadius: 20)
            .fill(type.color.secondary)
    }
}

// MARK: - CardType
extension ScreenCard {
    enum CardType: String, CaseIterable {
        case notes = "Notes"
        case learned = "Learned"
        case nonLearned = "Non-Learned"

        var color: (primary: Color, secondary: Color) {
            switch self {
            case .notes: return (.orangePrimary, .orangeSecondary)
            case .learned: return (.blueAccent, .blueSecondary)
            case .nonLearned: return (.purpleAccent, .purpleSecondary)
            }
        }

        var imageName: String {
            switch self {
            case .notes: return "book.closed.fill"
            case .learned: return "memories"
            case .nonLearned: return "pencil"
            }
        }
    }
}
