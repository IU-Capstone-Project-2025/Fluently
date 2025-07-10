//
//  StatisticInfoGrid.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI

final class StatisticInfoGrid: ObservableObject {
    var randomWeek: [Int] = []
    var randomMonth: [Int] = []
    var randomWYear: [Int] = []

    init() {
        generateRandoms()
    }

    func generateRandoms() {
        randomWeek = (0..<7).map { _ in Int.random(in: 1...20) }
        randomMonth = (0..<30).map { _ in Int.random(in: 1...20) }
        randomWYear = (0..<12).map { _ in Int.random(in: 10...50) }
    }

    func getAverage(range: TimeRange) -> Int {
        switch range {
            case .week:
                guard !randomWeek.isEmpty else { return 0 }
                let sum = randomWeek.reduce(0, +)
                return sum / randomWeek.count
            case .month:
                guard !randomMonth.isEmpty else { return 0 }
                let sum = randomMonth.reduce(0, +)
                return sum / randomMonth.count
            case .year:
                guard !randomWYear.isEmpty else { return 0 }
                let sum = randomWYear.reduce(0, +)
                return sum / randomWYear.count

        }
    }

    func getMax(range: TimeRange) -> Int {
        switch range {
            case .week : return randomWeek.max() ?? 1
            case .month: return randomMonth.max() ?? 1
            case .year: return randomWYear.max() ?? 1
        }
    }
}
