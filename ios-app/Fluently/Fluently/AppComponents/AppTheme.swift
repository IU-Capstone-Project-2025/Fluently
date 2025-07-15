//
//  AppTheme.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

// MARK: - App theme with colort
final class AppTheme {
//    Default
    let gray = Color("gray")
    let black = Color("blackFluently")
    let white = Color("whiteBackground")
    let textWhite = Color("whiteText")
    let textBlack = Color("whiteText")

//    Orange
    let orangePrimary = Color("orange.primary")
    let orangeSecondary = Color("orange.secondary")

//    Purple
    let purpleAccent = Color("purple.accent")
    let purpleSecondary = Color("purple.secondary")

//    Blue
    let blueAccent = Color("blue.accent")
    let blueSecondary = Color("blue.secondary")
}

// MARK: - Custom Fonts
final class AppFont {
    // Fraunces - primary
    let largeTitle = Font.custom("Fraunces", size: 34.0)
    let title = Font.custom("Fraunces", size: 28.0)
    let title2 = Font.custom("Fraunces", size: 22.0)
    let title3 = Font.custom("Fraunces", size: 20.0)

    let headline = Font.custom("Fraunces", size: 17.0)
    let callout = Font.custom("Fraunces", size: 16.0)
    let subheadline = Font.custom("Fraunces", size: 15.0)

    let body = Font.custom("Fraunces", size: 17.0)

    let footnote = Font.custom("Fraunces", size: 13.0)
    let caption = Font.custom("Fraunces", size: 12.0)
    let caption2 = Font.custom("Fraunces", size: 11.0)

    // Inter-Regular - secondary
    let secondaryLargeTitle = Font.custom("Inter-Regular", size: 34.0)
    let secondaryTitle = Font.custom("Inter-Regular", size: 28.0)
    let secondaryTitle2 = Font.custom("Inter-Regular", size: 22.0)
    let secondaryTitle3 = Font.custom("Inter-Regular", size: 20.0)

    let secondaryHeadline = Font.custom("Inter-Regular", size: 17.0)
    let secondaryCallout = Font.custom("Inter-Regular", size: 16.0)
    let secondarySubheadline = Font.custom("Inter-Regular", size: 15.0)

    let secondaryBody = Font.custom("Inter-Regular", size: 17.0)

    let secondaryFootnote = Font.custom("Inter-Regular", size: 13.0)
    let secondaryCaption = Font.custom("Inter-Regular", size: 12.0)
    let secondaryCaption2 = Font.custom("Inter-Regular", size: 11.0)
}

// MARK: - Extensions
extension Color {
    static let appTheme = AppTheme()
}

extension Font {
    static let appFont = AppFont()
}
