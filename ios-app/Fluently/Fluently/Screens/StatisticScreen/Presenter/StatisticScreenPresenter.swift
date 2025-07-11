//
//  StatisticScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 10.07.2025.
//

import Foundation
import SwiftUI

protocol StatisticScreenPresenting: ObservableObject {

}

final class StatisticScreenPresenter: StatisticScreenPresenting {
    @Published var selectedRange: TimeRange = .week
}
