//
//  AIChatScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation
import SwiftUI

final class AIChatScreenPresenter: ObservableObject{
    @Published var messages: [MessageModel] = MessageModel.mockGenerator()

    
}
