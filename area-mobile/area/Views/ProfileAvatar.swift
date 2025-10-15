//
//  ProfileAvatar.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct ProfileAvatar: View {
    let size: CGFloat
    let user: User?
    @State private var profileImage: UIImage?
    @StateObject private var authService = AuthService.shared
    
    init(size: CGFloat = 80, user: User?) {
        self.size = size
        self.user = user
    }
    
    var body: some View {
        ZStack {
            if let profileImage = profileImage {
                Image(uiImage: profileImage)
                    .resizable()
                    .aspectRatio(contentMode: .fill)
                    .frame(width: size, height: size)
                    .clipShape(Circle())
            } else {
                Circle()
                    .fill(OptimizedGradients.primaryGradient)
                    .frame(width: size, height: size)
                    .overlay(
                        Text((user ?? authService.currentUser)?.firstName?.prefix(1).uppercased() ?? "U")
                            .font(.system(size: size * 0.4, weight: .bold))
                            .foregroundColor(.white)
                    )
            }
        }
        .onAppear {
            loadProfileImage()
        }
        .onChange(of: user?.profileImage) { _, _ in
            loadProfileImage()
        }
        .onChange(of: authService.currentUser?.profileImage) { _, _ in
            loadProfileImage()
        }
    }
    
    private func loadProfileImage() {
        let currentUser = user ?? authService.currentUser
        guard let profileImagePath = currentUser?.profileImage else {
            profileImage = nil
            return
        }
        
        let fullURL = getFullImageURL(profileImagePath)
        if let url = fullURL {
            loadImageFromURL(url)
        }
    }
    
    private func loadImageFromURL(_ url: URL) {
        URLSession.shared.dataTask(with: url) { data, response, error in
            if let data = data, let image = UIImage(data: data) {
                DispatchQueue.main.async {
                    self.profileImage = image
                }
            }
        }.resume()
    }
    
    private func getFullImageURL(_ imagePath: String) -> URL? {
        if imagePath.hasPrefix("uploads/") {
            return URL(string: "\(AppConfig.baseURL)/\(imagePath)")
        }
        return URL(string: imagePath)
    }
}

#Preview {
    VStack(spacing: 20) {
        ProfileAvatar(size: 80, user: nil)
        ProfileAvatar(size: 60, user: nil)
        ProfileAvatar(size: 40, user: nil)
    }
    .background(Color.black)
}
