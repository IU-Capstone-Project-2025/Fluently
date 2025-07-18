//
//  DayWord.swift
//  Fluently
//
//  Created by Савва Пономарев on 14.07.2025.
//

import Foundation
import SwiftData

@Model
class DayWord {
    var date: Date
    var word: WordModel?

    init(word: WordModel?) {
        self.date = Date.now.startOfDay
        self.word = word
    }
}
