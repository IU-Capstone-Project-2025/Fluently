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

final class AppFont {
    let largeTitle = Font.custom("SedanSC-Regular", size: 34.0)
    let title = Font.custom("SedanSC-Regular", size: 28.0)
    let title2 = Font.custom("SedanSC-Regular", size: 22.0)
    let title3 = Font.custom("SedanSC-Regular", size: 20.0)

    let headline = Font.custom("SedanSC-Regular", size: 17.0)
    let callout = Font.custom("SedanSC-Regular", size: 16.0)
    let subheadline = Font.custom("SedanSC-Regular", size: 15.0)

    let footnote = Font.custom("SedanSC-Regular", size: 13.0)
    let caption = Font.custom("SedanSC-Regular", size: 12.0)
    let caption2 = Font.custom("SedanSC-Regular", size: 11.0)
}

extension Color {
    static let appTheme = AppTheme()
}

extension Font {
    static let appFont = AppFont()
}
