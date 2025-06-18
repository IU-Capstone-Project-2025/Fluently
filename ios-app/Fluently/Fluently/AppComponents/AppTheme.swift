//
//  AppTheme.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

struct AppTheme {
//    Default
    let gray = Color("gray")
    let black = Color("black")
    let white = Color("white")

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

extension Color {
    static let theme = AppTheme()
}
