//
//  TimeRange.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import Foundation

enum TimeRange: String, CaseIterable, Identifiable{
    case week = "Week"
    case month = "Month"
    case year = "Year"

    var id: String {self.rawValue}

    var number: Int {
        switch self {
            case .week: return 7
            case .month: return 30
            case .year: return 365
        }
    }
}
