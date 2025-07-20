//
//  CalendarScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import Foundation
import SwiftUI
import SwiftData

protocol CalendarScreenPresenting: ObservableObject {

}

final class CalendarScreenPresenter: CalendarScreenPresenting {
    @Published var selectedDate: Date

    var modelContext: ModelContext?

    init() {
        self.selectedDate = Date.now
    }

    func getForDate(_ date: Date) -> [WordModel] {
        guard let modelContext else {
            print("ModelContext not set")
            return []
        }

        let calendar = Calendar.current
        let startOfDay = calendar.startOfDay(for: date)
        guard let endOfDay = calendar.date(byAdding: .day, value: 1, to: startOfDay) else {
            return []
        }

        let predicate = #Predicate<WordModel> { word in
            return ((word.wordDate >= startOfDay) && (word.wordDate < endOfDay) && word.isInLibrary == true)
        }

        let descriptor = FetchDescriptor(
            predicate: predicate,
        )

        do {
            return try modelContext.fetch(descriptor)
        } catch {
            print("Failed to fetch words: \(error)")
            return []
        }
    }
    
    func setModelContext(_ context: ModelContext) {
        self.modelContext = context
    }
}
