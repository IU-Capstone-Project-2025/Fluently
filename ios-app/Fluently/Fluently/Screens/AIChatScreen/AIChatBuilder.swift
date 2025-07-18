//
//  AIChatBuilder.swift
//  Fluently
//
//  Created by Савва Пономарев on 17.07.2025.
//

import Foundation

enum AIChatBuilder {

    static func build (
        onExit: Optional< () -> Void>
    ) -> AIChatView {
        let interactor = AIChatInteractor()

        let presenter = AIChatScreenPresenter(
            interactor: interactor
        )

        return AIChatView (
            onExit: onExit,
            presenter: presenter
        )
    }
}
