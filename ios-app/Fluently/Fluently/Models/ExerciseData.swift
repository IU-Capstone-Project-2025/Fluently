//
//  ExerciseData.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.07.2025.
//

import Foundation
import SwiftData

// MARK: - Exercise Data Protocol
protocol ExerciseData: Codable {
    var correctAnswer: String { get }
}

struct EmptyExerciseData: ExerciseData {
    let correctAnswer: String = ""
}
