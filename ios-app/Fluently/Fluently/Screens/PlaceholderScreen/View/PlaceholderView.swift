//
//  PlaceholderView.swift
//  Fluently
//
//  Created by Савва Пономарев on 23.06.2025.
//

import Foundation
import SwiftUI

struct PlaceholderView: View {
    @Environment(\.dismiss) var dismiss
    var name: String?

    var body: some View {
        NavigationStack{
            VStack(alignment: .center) {
                Text(name ?? "Placeholder")
                    .font(.title)
            }
            .toolbar{
                ToolbarItem(placement: .topBarLeading) {
                    Image(systemName: "chevron.left")
                        .onTapGesture {
                            dismiss.callAsFunction()
                        }
                }
            }
        }
    }
}
