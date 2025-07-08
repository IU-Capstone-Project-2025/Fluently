//
//  MainView.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import Foundation
import SwiftUI

struct MainView: View {
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var accountData: AccountData

    @State var currentScreen: NavigationBar.Screen = .home

    var body: some View {
        ZStack {
            Group {
                switch currentScreen {
                    case .calendar:
                        CalendarScreenView()
                            .transition(
                                .asymmetric(
                                    insertion: .move(edge: .leading),
                                    removal: .move(edge: .leading)
                                )
                            )
                    case .home:
                        HomeScreenBuilder.build(router: router, acoount: accountData)
                            .transition(.opacity)
                    case .statistic:
                        StatisticScreenView()
                            .transition(
                                .asymmetric(
                                    insertion: .move(edge: .trailing),
                                    removal: .move(edge: .trailing)
                                )
                            )
                }
            }
            .animation(.easeInOut(duration: 0.3), value: currentScreen)
        }
        .toolbar {
            ToolbarItem(placement: .bottomBar) {
                NavigationBar(currentScreen: $currentScreen)
            }
        }
    }
}

struct MainView_Previews: PreviewProvider {
    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        @StateObject var account: AccountData = AccountData()
        @StateObject var router: AppRouter = AppRouter()


        var body: some View {
            MainView()
                .environmentObject(router)
                .environmentObject(account)
        }
    }
}
