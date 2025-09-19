//
//  HomeView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct HomeView: View {
    @State private var showTestView = false
    let onLogout: () -> Void
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                DecorativeShapes()
                
                VStack(spacing: 30) {
                    Spacer(minLength: 60)
                    
                    Text("AREA")
                        .font(AppTextStyles.title)
                        .foregroundColor(.white)
                        .padding(.bottom, 20)
                    
                    Text("Welcome to your AREA dashboard")
                        .font(AppTextStyles.subtitle)
                        .foregroundColor(.white)
                        .multilineTextAlignment(.center)
                        .padding(.horizontal, 24)
                    
                    VStack(spacing: 20) {
                        HomeActionButton(
                            title: "Test View",
                            icon: "testtube.2",
                            action: {
                                showTestView = true
                            }
                        )
                        
                        HomeActionButton(
                            title: "Settings",
                            icon: "gear",
                            action: {
                                print("Settings tapped")
                            }
                        )
                        
                        HomeActionButton(
                            title: "Profile",
                            icon: "person.circle",
                            action: {
                                print("Profile tapped")
                            }
                        )
                    }
                    .padding(.horizontal, 24)
                    
                    Spacer()
                    
                    Button(action: onLogout) {
                        HStack(spacing: 8) {
                            Image(systemName: "arrow.right.square")
                            Text("Logout")
                        }
                        .font(AppTextStyles.button)
                        .foregroundColor(.white)
                        .padding(.horizontal, 24)
                        .padding(.vertical, 12)
                        .background(
                            RoundedRectangle(cornerRadius: 8)
                                .fill(Color.red.opacity(0.8))
                        )
                    }
                    .padding(.bottom, 40)
                }
            }
        }
        .fullScreenCover(isPresented: $showTestView) {
            TestView()
        }
    }
}

struct HomeActionButton: View {
    let title: String
    let icon: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 16) {
                Image(systemName: icon)
                    .font(.system(size: 20))
                    .foregroundColor(.white)
                    .frame(width: 24)
                
                Text(title)
                    .font(AppTextStyles.button)
                    .foregroundColor(.white)
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
            }
            .padding(.horizontal, 20)
            .padding(.vertical, 16)
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(AppColors.darkBackground)
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                    )
            )
        }
    }
}

#Preview {
    HomeView(onLogout: {})
}
