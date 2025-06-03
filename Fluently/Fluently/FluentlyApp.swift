//
//  FluentlyApp.swift
//  Fluently
//
//  Created by Савва Пономарев on 03.06.2025.
//

import SwiftUI

@main
struct FluentlyApp: App {
    let persistenceController = PersistenceController.shared

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(\.managedObjectContext, persistenceController.container.viewContext)
        }
    }
}
