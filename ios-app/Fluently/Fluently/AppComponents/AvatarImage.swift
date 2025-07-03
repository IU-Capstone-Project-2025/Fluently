//
//  AvatarImage.swift
//  Fluently
//
//  Created by Савва Пономарев on 18.06.2025.
//

import SwiftUI

struct AvatarImage: View{
    // MARK: - Key Objects
    @EnvironmentObject var account: AccountData

    // MARK: - Properties
    let size: CGFloat
    let imageCache = NSCache<AnyObject, UIImage>()

    var onTap: (() -> Void)?

    var body: some View {
        Button {
            onTap?()
        } label: {
            if let imageUrlString = account.image,
               let imageUrl = URL(string: imageUrlString) {
                AsyncImage(url: imageUrl) { phase in
                    switch phase {
                        case .success(let image):
                            image
                                .resizable()
                                .scaledToFill()
                        case .failure(_):
                            fallbackIcon()
                        case .empty:
                            if let image = emptyImage(url: imageUrlString) {
                                image
                            } else {
                                ProgressView()
                            }
                        @unknown default:
                            fallbackIcon()
                    }
                }
            } else {
                fallbackIcon()
            }
        }
        .task {
            await cacheImage()
        }
        .clipShape(
            Circle()
        )
        .scaledToFit()
        .background(
            Circle()
                .fill(.orangeSecondary)
                .stroke(.orangeSecondary, lineWidth: 5)
                .frame(
                    width: size,
                    height: size
                )
        )
        .frame(width: size, height: size)
        .buttonStyle(.plain)
    }

    // image loading error handling
    private func fallbackIcon() -> some View {
        Image(systemName: "person")
            .resizable()
            .scaledToFit()
            .padding()
    }

    private func emptyImage(url imageUrlString: String) -> Image? {
        if let cache = imageCache.object(forKey: imageUrlString as AnyObject) {
            return Image(uiImage: cache)
        }
        return nil
    }

    private func cacheImage() async {
       guard let imageUrlString = account.image,
             let url = URL(string: imageUrlString) else { return }

       do {
           let (data, _) = try await URLSession.shared.data(from: url)
           if let image = UIImage(data: data) {
               await MainActor.run {
                   imageCache.setObject(image, forKey: imageUrlString as AnyObject)
               }
           }
       } catch {
           print("Failed to cache image: \(error.localizedDescription)")
       }
   }
}
