//
//  SheetViewModifier.swift
//  Fluently
//
//  Created by Савва Пономарев on 23.06.2025.
//

import SwiftUI

// Modifier for white sheet view 
struct SheetViewModifier: ViewModifier {
    var glassView: Bool = false

    // MARK: - View constances
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
            .glass(cornerRadius: Const.sheetCornerRadius)
            .background(
                UnevenRoundedRectangle(
                    topLeadingRadius: Const.sheetCornerRadius,
                    topTrailingRadius: Const.sheetCornerRadius,
                )
                .fill(
                    glassView ? .clear :
                    .whiteBackground
                )
            )
            .ignoresSafeArea(.all)
    }
}

struct SheetViewPreview: PreviewProvider {

    static var previews: some View {
        VStack {
            Spacer()
                .frame(
                    height: 200
                )
            Text("preview")
                .modifier(SheetViewModifier(glassView: true))
        }
            .modifier(BackgroundViewModifier(colorful: true))
    }
}
