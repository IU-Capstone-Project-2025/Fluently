//
//  NotesScreenPresenter.swift
//  Fluently
//
//  Created by Савва Пономарев on 01.07.2025.
//

import Foundation
import SwiftUI

protocol NotesScreenPresenting: ObservableObject {

    func toggleRecording()
}

final class NotesScreenPresenter: NotesScreenPresenting {
    let router: NotesScreenRouter

    @Published var isRecording = false

    init(router: NotesScreenRouter) {
        self.router = router
    }

    func toggleRecording() {
        switch isRecording {
            case true:
                stopRecording()
            case false:
                startRecording()
        }
    }

}

// MARK: - Private
private extension NotesScreenPresenter {
    func startRecording() {
        isRecording = true

        print("Start")
    }

    func stopRecording() {
        isRecording = false

        print("Stop")
    }
}
