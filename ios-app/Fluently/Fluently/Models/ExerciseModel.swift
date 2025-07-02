//
//  ExerciseModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

class ExerciseModel: Codable{
    var data: String
    var type: String

    var correctAnswer: String?

    init(
        data: String,
        type: String
    ) {
        self.data = data
        self.type = type
    }
}

