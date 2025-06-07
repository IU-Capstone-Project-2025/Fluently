//
//  AppTheme.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import Foundation
import SwiftUI

struct AppTheme {
    let primary = Color("greenColor")
    let secondary = Color("purpleColor")
    let tertiary = Color("redColor")
    let complementary1 = Color("blueColor")
    let complementary2 = Color("yellowColor")
    let backgroundColor = Color("background")
}

extension Color {
    static let theme = AppTheme()
}
