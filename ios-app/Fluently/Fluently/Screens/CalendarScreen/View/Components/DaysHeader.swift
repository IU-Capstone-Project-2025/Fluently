//
//  DaysHeader.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import SwiftUI

struct DaysHeader: View {
    @Binding var selectedDate: Date

    @State private var offset: CGFloat = 0
    @State private var isSwiping = false

    @State var weekStart: Date = Date.now.startOfWeek!

    var dayNameFormatter: DateFormatter {
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "e"
        return dateFormatter
    }

    var dayNumberFormatter: DateFormatter{
        let dateFormatter = DateFormatter()
        dateFormatter.dateFormat = "d"
        return dateFormatter
    }

    var body: some View {
        weekHStack()
            .offset(x: offset)
            .animation(.interactiveSpring(), value: offset)
            .gesture(
                DragGesture()
                    .onChanged { gesture in
                        isSwiping = true
                        offset = gesture.translation.width
                    }
                    .onEnded { gesture in
                        if abs(gesture.translation.width) > 100 {
                            if gesture.translation.width > 0 {
                                weekStart = Calendar.current.date(byAdding: .weekOfYear, value: -1, to: weekStart) ?? selectedDate
                            } else {
                                weekStart = Calendar.current.date(byAdding: .weekOfYear, value: 1, to: weekStart) ?? selectedDate
                            }
                        }
                        offset = 0
                        isSwiping = false
                    }
            )
    }

    func weekHStack() -> some View {
        HStack(spacing: 16) {
            ForEach(0..<7) { dayOffset in
                if let dayDate = Calendar.current.date(byAdding: .day, value: dayOffset, to: weekStart) {
                    dayItem(date: dayDate) {
                        selectedDate = dayDate
                    }
                }
            }
        }
        .onAppear{
            selectedDate = Date.now
        }
        .padding()
    }

    @ViewBuilder
    func dayItem(date: Date, onSelect: @escaping () -> Void) -> some View {
        let todayStart = Calendar.current.startOfDay(for: Date.now)
        let dateStart = Calendar.current.startOfDay(for: date)
        let selectedStart = Calendar.current.startOfDay(for: selectedDate)

        let isToday = todayStart == dateStart
        let isSelected = selectedStart == dateStart

        var numberColor: Color {
            if isToday {
                return isSelected ? .blackText : .pink
            }
            return .blackText
        }

        var backgroundDate: Color {
            if isSelected || (selectedStart == todayStart && isToday) {
                return .purpleAccent
            }
            return .white
        }

        var shouldShowCircle: Bool {
            isSelected || (selectedStart == todayStart && isToday)
        }

        VStack(spacing: 8) {
            Text(date.formatted(Date.FormatStyle().weekday(.narrow)))
                .foregroundStyle(.blackText)
            Text(dayNumberFormatter.string(from: date))
                .foregroundStyle(numberColor)
                .frame(
                    maxWidth: .infinity,
                    maxHeight: 40
                )
                .glass(
                    cornerRadius: 100,
                    fill: shouldShowCircle ? Color.purpleAccent : Color.clear
                )
        }
        .onTapGesture {
            selectedDate = date
        }
    }
}
