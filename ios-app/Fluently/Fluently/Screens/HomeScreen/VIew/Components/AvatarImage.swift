//
//  AvatarImage.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI

struct AvatarImage: View{
    let size: CGFloat
    @State var icon: Image = Image(systemName: "person")

    var body: some View {
        Button {
//            TODO: open profile
        } label: {
            icon
                .resizable()
                .scaledToFit()
                .padding()
                .clipShape(
                    Circle()
                )
                .background(
                    Circle()
                        .fill(.orangeSecondary)
                )
                .frame(width: size, height: size)
        }
        .buttonStyle(.plain)
    }
}
