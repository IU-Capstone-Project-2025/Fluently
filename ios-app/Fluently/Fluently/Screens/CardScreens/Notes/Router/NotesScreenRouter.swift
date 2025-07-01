//
//  NotesScreenRouter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

final class NotesScreenRouter {
    @ObservedObject var router: AppRouter

    init(router: AppRouter) {
        self.router = router
    }
}
