//
//  MainView.swift
//  Fluently
//
//  Created by Савва Пономарев on 08.07.2025.
//

import Foundation
import SwiftUI

// MARK: - Main Screen
struct MainView: View {
    @EnvironmentObject var router: AppRouter
    @EnvironmentObject var accountData: AccountData

    @State var currentScreen: NavigationBar.Screen = .home

    @State private var calendarView: CalendarScreenView?
    @State private var homeView: HomeScreenView?
    @State private var statisticView:StatisticScreenView?

    /// setup screens view
    ///  `CalendarView` -  screen with info day by day
    ///  `HomeView` - screen with main functions
    ///  `StatisticView` - screen with user statistic
    func setupScreens() {
        calendarView = CalendarScreenBuilder.build()
        homeView = HomeScreenBuilder.build(router: router, acoount: accountData)
        statisticView = StatisticScreenBuilder.build()
    }

    var body: some View {
        ZStack {
            Group {
                switch currentScreen {
                    case .calendar:
                        calendarView
                            .transition(
                                .asymmetric(
                                    insertion: .move(edge: .leading),
                                    removal: .move(edge: .leading)
                                )
                            )
                    case .home:
                        homeView
                            .transition(.opacity)
                    case .statistic:
                        statisticView
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
        .onAppear {
            setupScreens()
        }
        .toolbar {
            ToolbarItem(placement: .bottomBar) {
                NavigationBar(currentScreen: $currentScreen)
            }
        }
    }
}

// MARK: - Preview Provider
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
