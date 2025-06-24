//
//  BackgroundViewModifier.swift
//  Fluently
//
//  Created by Савва Пономарев on 24.06.2025.
//

import Foundation
import SwiftUI

struct BackgroundViewModifier: ViewModifier {
    func body(content: Content) -> some View {
        content
            .containerRelativeFrame([.horizontal, .vertical])
            .background(.orangePrimary)
    }
}
