//
//  StatisticInfoGrid.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI
import SwiftData

final class StatisticInfoGrid: ObservableObject {
    // MARK: - Properties

    @Published var weekStat: [Int] = Array(repeating: 0, count: 7)
    @Published var monthStat: [Int] = Array(repeating: 0, count: 31)
    @Published var yearStat: [Int] = Array(repeating: 0, count: 12)

    // MARK: - Components for data
    var modelContext: ModelContext?
    private let calendar = Calendar.current

    init() {
#if targetEnvironment(simulator)
        generateRandoms()
#endif
    }

    // MARK: - Main stat info

    /// calculate average number of words on differen time ranges
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

    /// get max number of words on time range
    func getMax(range: TimeRange) -> Int {
        switch range {
            case .week : return weekStat.max() ?? 1
            case .month: return monthStat.max() ?? 1
            case .year: return yearStat.max() ?? 1
        }
    }

    
    // MARK: - Statistic calculation

    /// Calculate statistic for the week
    func loadWeekStatistic() {
        guard let context = modelContext else { return }

        guard let weekStart = Date.now.startOfWeek else { return }

        guard let weekEnd = calendar.date(byAdding: .day, value: 6, to: weekStart) else { return }
        let endOfWeek = calendar.date(bySettingHour: 23, minute: 59, second: 59, of: weekEnd)!

        let predicate = #Predicate<WordModel> { word in
            return word.wordDate >= weekStart && word.wordDate <= endOfWeek && word.isLearned
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

    /// Calculate statistic for the month
    func loadMonthStatistic() {
        guard let context = modelContext else { return }

        let monthStart = Date.now.startOfMonth
        let monthEnd = Date.now.endOfMonth

        let predicate = #Predicate<WordModel> { word in
            return word.wordDate >= monthStart && word.wordDate <= monthEnd && word.isLearned
        }

        do {
            let words = try context.fetch(FetchDescriptor(predicate: predicate))

            let daysInMonth = calendar.range(of: .day, in: .month, for: Date.now)?.count ?? 31
            var dailyCounts = Array(repeating: 0, count: daysInMonth)

            for word in words {
                let learnedDate = word.wordDate
                let day = calendar.component(.day, from: learnedDate)
                let index = day - 1
                dailyCounts[index] += 1
            }

            monthStat = dailyCounts
        } catch {
            print("Failed to fetch month statistics: \(error)")
            monthStat = Array(repeating: 0, count: Date.now.getLastDayOfMonth())
        }
    }

    /// Calculate statistic for the year
    func loadYearStatistic() {
        guard let context = modelContext else { return }

        let yearStart = Date.now.startOfYear
        let yearEnd = Date.now.endOfYear

        let predicate = #Predicate<WordModel> { word in
            return word.wordDate >= yearStart && word.wordDate <= yearEnd && word.isLearned
        }

        do {
            let words = try context.fetch(FetchDescriptor(predicate: predicate))

            var monthlyCounts = Array(repeating: 0, count: 12)

            for word in words {
                let learnedDate = word.wordDate
                let month = calendar.component(.month, from: learnedDate)
                let index = month - 1
                monthlyCounts[index] += 1
            }

            yearStat = monthlyCounts
        } catch {
            print("Failed to fetch year statistics: \(error)")
            yearStat = Array(repeating: 0, count: 12)
        }
    }

    /// setup statistic
    func setModelContext(_ context: ModelContext) {
        self.modelContext = context

        loadWeekStatistic()
        loadMonthStatistic()
        loadYearStatistic()
    }

    // MARK: - Mock data
    func generateRandoms() {
        weekStat = (0..<7).map { _ in Int.random(in: 1...20) }
        monthStat = (0..<30).map { _ in Int.random(in: 1...20) }
        yearStat = (0..<12).map { _ in Int.random(in: 10...50) }
    }
}
