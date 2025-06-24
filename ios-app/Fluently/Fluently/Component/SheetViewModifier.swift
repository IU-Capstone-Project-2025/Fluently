//
//  SheetViewModifier.swift
//  Fluently
//
//  Created by Савва Пономарев on 23.06.2025.
//

import SwiftUI

struct SheetViewModifier: ViewModifier {

    private enum Const {
        // Paddings
        static let horizontalPadding = CGFloat(30)

        // Corner Radiuses
        static let sheetCornerRadius = CGFloat(20)
        static let gridInfoVerticalPadding = CGFloat(20)
    }


    func body(content: Content) -> some View {
        content
            .padding(.top, Const.gridInfoVerticalPadding)
            .frame(maxWidth: .infinity, maxHeight: .infinity)
            .background(
                UnevenRoundedRectangle(
                    topLeadingRadius: Const.sheetCornerRadius,
                    topTrailingRadius: Const.sheetCornerRadius,
                )
                .fill(
                    .whiteBackground
                )
                .ignoresSafeArea(.all)
            )
    }
}
