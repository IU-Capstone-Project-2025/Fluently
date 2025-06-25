//
//  AvatarImage.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI

struct AvatarImage: View{
    // MARK: - Key Objects
    @EnvironmentObject var account: AccountData

    // MARK: - Properties
    let size: CGFloat

    var onTap: (() -> Void)?

    var body: some View {
        Button {
            onTap?()
        } label: {
            if let imageUrlString = account.image,
               let imageUrl = URL(string: imageUrlString) {
                AsyncImage(url: imageUrl) { phase in
                    switch phase {
                    case .success(let image):
                        image
                            .resizable()
                            .scaledToFill()
                    case .failure(_):
                        fallbackIcon()
                    case .empty:
                        ProgressView()
                    @unknown default:
                        fallbackIcon()
                    }
                }
            } else {
                fallbackIcon()
            }
        }
        .clipShape(
            Circle()
        )
        .scaledToFit()
        .padding(3)
        .background(
            Circle()
                .fill(.orangeSecondary)
        )
        .frame(width: size, height: size)
        .buttonStyle(.plain)
    }

    // image loading error handling
    private func fallbackIcon() -> some View {
        Image(systemName: "person")
            .resizable()
            .scaledToFit()
            .padding()
    }
}
