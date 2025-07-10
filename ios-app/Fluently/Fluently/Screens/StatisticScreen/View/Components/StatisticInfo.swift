//
//  StatisticInfo.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import SwiftUI

struct StatisticInfo: View {
    var range: TimeRange

    @StateObject var presenter = StatisticInfoGrid()

    @State var weekStart = Date.now.startOfWeek

    var dayNumberFormatter: DateFormatter {
        let formatter = DateFormatter()
        formatter.dateFormat = "d"
        return formatter
    }

    var monthFormatter: DateFormatter {
        let formatter = DateFormatter()
        formatter.dateFormat = "MMMM"
        return formatter
    }

    var yearFormatter: DateFormatter {
        let formatter = DateFormatter()
        formatter.dateFormat = "yyyy"
        return formatter
    }


    var body: some View {
        dayHeader
        infoGrid
    }

    var dayHeader: some View {
        HStack {
            VStack(alignment: .leading) {
                Text("average")
                    .font(.appFont.secondarySubheadline)
                    .foregroundStyle(.blackText)
                HStack(alignment: .lastTextBaseline) {
                    Text("\(presenter.getAverage(range: range))")
                        .font(.appFont.largeTitle.bold())
                        .foregroundStyle(.orangePrimary)
                    Text("words")
                        .font(.appFont.secondarySubheadline)
                        .foregroundStyle(.blackText)
                }
                rangeLabel
                    .font(.appFont.secondarySubheadline)
                    .foregroundStyle(.blackText)
            }
            .padding()
            .glass(cornerRadius: 20, fill: .orangePrimary)
        }
        .padding(.horizontal)
        .frame(
            maxWidth: .infinity,
            alignment: .leading,
        )
    }

    @ViewBuilder
    var rangeLabel: some View {
        switch range {
            case .week:
                Text(
                    "\(dayNumberFormatter.string(from: weekStart!)) – \(dayNumberFormatter.string(from: weekStart!.endOfWeek!))")
                Text("\(monthFormatter.string(from: weekStart!).lowercased()) \(yearFormatter.string(from: weekStart!).lowercased())")
            case .month:
                Text("\(monthFormatter.string(from: weekStart!).lowercased()) \(yearFormatter.string(from: weekStart!).lowercased())")
            case .year:
                Text("\(yearFormatter.string(from: weekStart!).lowercased())")
        }
    }

    var infoGrid: some View {
        ZStack {
            VStack{
                rangeGrid
                rangeFooter
            }
            .padding(8)
        }
        .frame(
            maxWidth: .infinity,
            maxHeight: .infinity
        )
        .glass(
            cornerRadius: 20,
            fill: .orangeSecondary
        )
        .padding()
    }

    var dayNameFormatter: DateFormatter {
        let formatter = DateFormatter()
        formatter.dateFormat = "EE"
        return formatter
    }

    @ViewBuilder
    var rangeFooter: some View {
        HStack (alignment: .center) {
            switch range {
                case .week:
                    ForEach(0..<7) { dayOffset in
                        if let dayDate = Calendar.current.date(byAdding: .day, value: dayOffset, to: weekStart!) {
                            Text(dayNameFormatter.string(from: dayDate))
                                .frame(maxWidth: .infinity)
                        }
                    }
                case .month:
                    GeometryReader { geometry in
                        let daysInMonth = weekStart!.getLastDayOfMonth()
                        let cellWidth = geometry.size.width / 30

                        ForEach(0..<30) { dayOffset in
                            if dayOffset % 7 == 0 && 1 + dayOffset < weekStart!.getLastDayOfMonth() {
                                Text("\(1 + dayOffset)")
                                    .frame(alignment: .leading)
                                    .offset(x: CGFloat(dayOffset % 30) * cellWidth)
                            }
                        }
                    }
                    .frame(height: 20)
                case .year:
                    ForEach(0..<12) { monthOffset in
                        Text(Calendar.current.monthSymbols[monthOffset].prefix(1))
                            .frame(maxWidth: .infinity)
                    }
            }
        }
        .bold()
        .foregroundStyle(.orangePrimary)
    }

    @ViewBuilder
    var rangeGrid: some View {
        let maxValue = presenter.getMax(range: range)

        HStack {
            switch range {
                case .week:
                    ForEach(0..<7, id: \.self) { index in
                        let value = presenter.randomWeek[index]
                        VStack {
                            RoundedRectangle(cornerRadius: 4)
                                .fill(Color.orangePrimary)
                                .frame(height: CGFloat(value) / CGFloat(maxValue) * 100)
                                .frame(maxHeight: 120)
                        }
                        .background(Color.orangePrimary.opacity(0.2))
                        .cornerRadius(4)
                    }

                case .month:
                    ForEach(0..<30, id: \.self) { index in
                        let value = presenter.randomMonth[index]
                        VStack {
                            RoundedRectangle(cornerRadius: 4)
                                .fill(Color.orangePrimary)
                                .frame(height: CGFloat(value) / CGFloat(maxValue) * 100)
                                .frame(maxHeight: 120)
                        }
                        .background(index % 7 == 0 ? Color.orangePrimary.opacity(0.3) : Color.clear)
                        .cornerRadius(4)
                    }
                    

                case .year:
                    ForEach(0..<12, id: \.self) { index in
                        let value = presenter.randomWYear[index]
                        VStack {
                            RoundedRectangle(cornerRadius: 4)
                                .fill(Color.orangePrimary)
                                .frame(height: CGFloat(value) / CGFloat(maxValue) * 100)
                                .frame(maxHeight: 120)
                        }
                        .background(Color.orangePrimary.opacity(0.2))
                        .cornerRadius(4)
                    }
            }
        }
    }
}


extension Date {
    func getLastDayOfMonth() -> Int {
        let calendar = Calendar.current
        let range = calendar.range(of: .day, in: .month, for: self)!
        return range.upperBound - 1
    }
}
