//
//  DateExt.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import Foundation

// MARK: - Date extension
extension Date {
    var startOfWeek: Date? {
        let calendar = Calendar.current
        guard let day = calendar.date(from: calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: self)) else { return nil }
//        return calendar.date(byAdding: .day, value: 1, to: sunday)
        return day
    }

    var endOfWeek: Date? {
        let calendar = Calendar.current
        guard let day = calendar.date(from: calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: self)) else { return nil }
//        return calendar.date(byAdding: .day, value: 7, to: day)
        return day
    }

    func addingWeeks(_ weeks: Int) -> Date {
        Calendar.current.date(byAdding: .weekOfYear, value: weeks, to: self) ?? self
    }
}

enum Weekday: Int {
    case mon = 1
    case tue = 2
    case wed = 3
    case thu = 4
    case fri = 5
    case sat = 6
    case sun = 7

    var shortName: String {
       return Calendar.current.shortWeekdaySymbols[self.rawValue - 1]
   }
}
