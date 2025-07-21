//
//  NetworkSheetView.swift
//  Fluently
//
//  Created by Савва Пономарев on 21.07.2025.
//

import Foundation
import SwiftUI

struct NetworkSheetView: View {

    @Environment(\.isNetworkConnected) var isConnected
    @Environment(\.connectionType) var connectionType

    var body: some View {
        VStack {
            Spacer()
            VStack(spacing: 20) {
                Image(systemName: "wifi.exclamationmark")
                    .font(.system(size: 80, weight: .semibold))
                    .frame(height: 100)

                Text("No Internet Connectivity")
                    .font(.appFont.title3)
                    .fontWeight(.semibold)

                Text("Please check your internet connection \nto continue using the app")
                    .multilineTextAlignment(.center)
                    .font(.appFont.headline)
            }
            .padding(20)
            .glass(
                cornerRadius: 20,
                fill: .orangePrimary,
                opacity: 0.8
            )
            .background(
                RoundedRectangle(cornerRadius: 20)
                    .fill(.ultraThinMaterial)
            )
        }
    }
}

// MARK: - Preview Provider
struct NetworkSheetPreview: PreviewProvider {

    static var previews: some View {
        PreviewWrapper()
    }

    struct PreviewWrapper: View {
        @State var bebe = true

        var body: some View {
            VStack {
                Circle()
                    .fill(.red)
                Text("Text")
                    .onTapGesture {
                        bebe = true
                    }
            }
            .sheet(isPresented: $bebe) {
                NetworkSheetView()
                    .presentationDetents([.medium])
                    .presentationBackground(.clear)
                    .presentationBackgroundInteraction(.disabled)
                    .interactiveDismissDisabled()
            }

//            .fullScreenCover(isPresented: $bebe) {
//                NetworkSheetView()
//                    .presentationDetents([.medium])
//                    .presentationBackground(.clear)
//                    .presentationBackgroundInteraction(.disabled)
//                    .interactiveDismissDisabled()
//            }
        }
    }
}
