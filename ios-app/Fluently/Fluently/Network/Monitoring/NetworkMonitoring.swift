//
//  NetworkMonitoring.swift
//  Fluently
//
//  Created by Савва Пономарев on 20.07.2025.
//

import SwiftUI
import Network

extension EnvironmentValues {
    @Entry var isNetworkConnected: Bool?
    @Entry var connectionType: NWInterface.InterfaceType?
}

class NetworkMonitoring: ObservableObject {

    @Published var isNetworkConnected: Bool?
    @Published var connectionType: NWInterface.InterfaceType?

    private var queue = DispatchQueue(label: "Monitoring")
    private var monitor = NWPathMonitor()

    init() {
        startMonitoring()
    }

    private func startMonitoring() {
        monitor.pathUpdateHandler = { path in
            Task { @MainActor in
                self.isNetworkConnected = path.status == .satisfied

                let types: [NWInterface.InterfaceType] = [.wifi, .cellular, .loopback, .wiredEthernet]
                if let type = types.first(where: { path.usesInterfaceType($0) }) {
                    self.connectionType = type
                } else {
                    self.connectionType = nil
                }
            }
        }

        monitor.start(queue: queue)
    }

    func stopMonitoring() {
         
    }
}
