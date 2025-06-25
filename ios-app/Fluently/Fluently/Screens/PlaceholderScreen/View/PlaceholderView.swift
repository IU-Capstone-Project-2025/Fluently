//
//  PlaceholderView.swift
//  Fluently
//
//  Created by Савва Пономарев on 23.06.2025.
//

import Foundation
import SwiftUI

struct PlaceholderView: View {
    var name: String?

    var body: some View {
        Text(name ?? "Placeholder")
            .font(.title)
    }
}
