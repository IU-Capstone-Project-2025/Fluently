//
//  StatisticInfoGrid.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI
import SwiftData

final class StatisticInfoGrid: ObservableObject {
    @Published var weekStat: [Int] = Array(repeating: 0, count: 7)
    @Published var monthStat: [Int] = Array(repeating: 0, count: 31)
    @Published var yearStat: [Int] = Array(repeating: 0, count: 12)

    var modelContext: ModelContext?
    private let calendar = Calendar.current

    init() {
#if targetEnvironment(simulator)
        generateRandoms()
#endif
    }

    func loadWeekStatistics(weekStart: Date) {
        guard let context = modelContext else { return }
        
        guard let weekEnd = calendar.date(byAdding: .day, value: 6, to: weekStart) else { return }
        let endOfWeek = calendar.date(bySettingHour: 23, minute: 59, second: 59, of: weekEnd)!

        let predicate = #Predicate<WordModel> { word in
            return word.wordDate >= weekStart && word.wordDate <= endOfWeek
        }

        do {
            let words = try context.fetch(FetchDescriptor(predicate: predicate))

            var dailyCounts = Array(repeating: 0, count: 7)

            for word in words {
                let learnedDate = word.wordDate
                let components = calendar.dateComponents([.weekday], from: learnedDate)
                if let weekday = components.weekday {
                    let index = (weekday - calendar.firstWeekday + 7) % 7
                    dailyCounts[index] += 1
                }
            }

            weekStat = dailyCounts
        } catch {
            print("Failed to fetch week statistics: \(error)")
            weekStat = Array(repeating: 0, count: 7)
        }
    }

    func setModelContext(_ context: ModelContext) {
        self.modelContext = context

        loadWeekStatistics(weekStart: Date.now.startOfWeek!)
    }

    func generateRandoms() {
        weekStat = (0..<7).map { _ in Int.random(in: 1...20) }
        monthStat = (0..<30).map { _ in Int.random(in: 1...20) }
        yearStat = (0..<12).map { _ in Int.random(in: 10...50) }
    }

    func getAverage(range: TimeRange) -> Int {
        switch range {
            case .week:
                guard !weekStat.isEmpty else { return 0 }
                let sum = weekStat.reduce(0, +)
                return sum / weekStat.count
            case .month:
                guard !monthStat.isEmpty else { return 0 }
                let sum = monthStat.reduce(0, +)
                return sum / monthStat.count
            case .year:
                guard !yearStat.isEmpty else { return 0 }
                let sum = yearStat.reduce(0, +)
                return sum / yearStat.count

        }
    }

    func getMax(range: TimeRange) -> Int {
        switch range {
            case .week : return weekStat.max() ?? 1
            case .month: return monthStat.max() ?? 1
            case .year: return yearStat.max() ?? 1
        }
    }
}
