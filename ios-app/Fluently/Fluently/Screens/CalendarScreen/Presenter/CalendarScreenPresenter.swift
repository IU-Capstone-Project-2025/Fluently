//
//  CalendarScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 09.07.2025.
//

import Foundation
import SwiftUI

protocol CalendarScreenPresenting: ObservableObject {

}

final class CalendarScreenPresenter: CalendarScreenPresenting {
    @Published var selectedDate: Date

    init() {
        self.selectedDate = Date.now
    }
}
