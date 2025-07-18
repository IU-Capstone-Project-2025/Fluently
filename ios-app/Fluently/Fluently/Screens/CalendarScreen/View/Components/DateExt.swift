//
//  DateExt.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import Foundation

// MARK: - Date extension
extension Date {
    // MARK: - Week Helpers
    var startOfWeek: Date? {
        let calendar = Calendar.current
        guard let day = calendar.date(from: calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: self)) else { return nil }
        return day
    }

    var endOfWeek: Date? {
        let calendar = Calendar.current
        guard let day = calendar.date(from: calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: self)) else { return nil }
        return calendar.date(byAdding: .day, value: 6, to: day)
    }

    func addingWeeks(_ weeks: Int) -> Date {
        Calendar.current.date(byAdding: .weekOfYear, value: weeks, to: self) ?? self
    }

    // MARK: - Month Helpers
    var startOfMonth: Date {
        Calendar.current.date(
            from: Calendar.current.dateComponents([.year, .month], from: self)
        ) ?? self
    }

    var endOfMonth: Date {
        let calendar = Calendar.current
        guard let nextMonth = calendar.date(byAdding: .month, value: 1, to: self.startOfMonth) else { return self }
        return calendar.date(byAdding: .day, value: -1, to: nextMonth) ?? self
    }

    var startOfDay: Date {
        Calendar.current.startOfDay(for: self)
    }

    var endOfDay: Date {
        var components = DateComponents()
        components.day = 1
        components.second = -1
        return Calendar.current.date(byAdding: components, to: startOfDay) ?? self
    }

    var startOfYear: Date {
        let calendar = Calendar.current
        let components = calendar.dateComponents([.year], from: self)
        return calendar.date(from: components)!
    }

    var endOfYear: Date {
        let calendar = Calendar.current
        var components = DateComponents()
        components.year = 1
        components.second = -1
        return calendar.date(byAdding: components, to: startOfYear)!
    }
}
