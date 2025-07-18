//
//  RangeHeader.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI

struct RangeHeader: View {
    @Binding var selectedRange: TimeRange

    var body: some View {
        weekHStack()
    }

    func weekHStack() -> some View {
        HStack(alignment: .top, spacing: 16) {
            Picker("Select range", selection: $selectedRange) {
                ForEach(TimeRange.allCases, id: \.id) { range in
                    Text(range.rawValue).tag(range)
                        .foregroundStyle(.blackText)
                }
            }
            .pickerStyle(.segmented)
        }
        .padding()
    }
}



