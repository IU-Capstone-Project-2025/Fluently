//
//  ExerciseModel.swift
//  Fluently
//
//  Created by Савва Пономарев on 02.07.2025.
//

import Foundation

final class ExerciseModel: Codable{
    var data: String
    var type: String

    init(
        data: String,
        type: String
    ) {
        self.data = data
        self.type = type
    }
}

