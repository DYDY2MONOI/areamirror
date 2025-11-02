//
//  ProfileView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct ProfileView: View {
    @StateObject private var authService = AuthService.shared
    @State private var showLoginView = false
    @State private var showRegisterView = false
    @State private var showEditProfile = false
    @State private var showNotifications = false
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                if authService.isAuthenticated {
                    ScrollView {
                        VStack(spacing: 24) {
                            VStack(spacing: 16) {
                                ProfileAvatar(size: 80, user: authService.currentUser)
                                
                                VStack(spacing: 4) {
                                    Text("\(authService.currentUser?.firstName ?? "") \(authService.currentUser?.lastName ?? "")")
                                        .font(.title2)
                                        .fontWeight(.bold)
                                        .foregroundColor(.white)
                                    
                                    Text(authService.currentUser?.email ?? "")
                                        .font(.body)
                                        .foregroundColor(.gray)
                                }
                            }
                            .padding(.top, 20)
                            
                            VStack(spacing: 12) {
                                ProfileMenuItem(
                                    icon: "person.circle",
                                    title: "Edit Profile",
                                    action: {
                                        showEditProfile = true
                                    }
                                )
                                
                                ProfileMenuItem(
                                    icon: "bell",
                                    title: "Notifications",
                                    action: {
                                        showNotifications = true
                                    }
                                )
                                
                                ProfileMenuItem(
                                    icon: "gear",
                                    title: "Settings",
                                    action: {
                                    }
                                )
                                
                                ProfileMenuItem(
                                    icon: "questionmark.circle",
                                    title: "Help & Support",
                                    action: {
                                    }
                                )
                                
                                ProfileMenuItem(
                                    icon: "info.circle",
                                    title: "About",
                                    action: {
                                    }
                                )
                            }
                            .padding(.horizontal, 20)
                            
                            Button(action: {
                                authService.logout()
                            }) {
                                HStack {
                                    Image(systemName: "power")
                                        .font(.title3)
                                    
                                    Text("LOGOUT")
                                        .font(.headline)
                                        .fontWeight(.semibold)
                                }
                                .foregroundColor(.red)
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(
                                    RoundedRectangle(cornerRadius: 12)
                                        .fill(Color.red.opacity(0.1))
                                        .overlay(
                                            RoundedRectangle(cornerRadius: 12)
                                                .stroke(Color.red.opacity(0.3), lineWidth: 1)
                                        )
                                )
                            }
                            .padding(.horizontal, 20)
                            .padding(.top, 20)
                        }
                        .padding(.bottom, 40)
                    }
                } else {
                    VStack(spacing: 30) {
                        Spacer()
                        VStack(spacing: 16) {
                            Text("AREA")
                                .font(.system(size: 48, weight: .bold, design: .rounded))
                                .foregroundColor(.white)
                            
                            Text("Automate your workflow")
                                .font(.title3)
                                .foregroundColor(.gray)
                        }
                        
                        VStack(spacing: 16) {
                            Button(action: {
                                showLoginView = true
                            }) {
                                Text("LOGIN")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .padding(.vertical, 16)
                                    .background(AppGradients.button)
                                    .cornerRadius(12)
                            }
                            
                            Button(action: {
                                showRegisterView = true
                            }) {
                                Text("SIGN UP")
                                    .font(.headline)
                                    .foregroundColor(AppColors.primaryBlue)
                                    .frame(maxWidth: .infinity)
                                    .padding(.vertical, 16)
                                    .background(
                                        RoundedRectangle(cornerRadius: 12)
                                            .stroke(AppColors.primaryBlue, lineWidth: 2)
                                    )
                            }
                        }
                        .padding(.horizontal, 40)
                        
                        Spacer()
                    }
                }
            }
        }
        .navigationTitle("Profile")
        .navigationBarHidden(true)
        .sheet(isPresented: $showLoginView) {
            LoginView(
                onLoginSuccess: {
                    showLoginView = false
                },
                onSignUpTap: {
                    showLoginView = false
                    showRegisterView = true
                }
            )
        }
        .sheet(isPresented: $showRegisterView) {
            RegisterView(
                onRegisterSuccess: {
                    showRegisterView = false
                },
                onLoginTap: {
                    showRegisterView = false
                    showLoginView = true
                }
            )
        }
        .sheet(isPresented: $showEditProfile) {
            EditProfileView(
                onDismiss: {
                    showEditProfile = false
                }
            )
        }
        .sheet(isPresented: $showNotifications) {
            NotificationsView()
        }
        .onAppear {
            if authService.isAuthenticated {
                authService.fetchProfile()
            }
        }
    }
}

struct ProfileMenuItem: View {
    let icon: String
    let title: String
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            HStack(spacing: 16) {
                Image(systemName: icon)
                    .font(.title3)
                    .foregroundColor(AppColors.primaryBlue)
                    .frame(width: 24)
                
                Text(title)
                    .font(.body)
                    .foregroundColor(.white)
                
                Spacer()
                
                Image(systemName: "chevron.right")
                    .font(.caption)
                    .foregroundColor(.gray)
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(AppColors.darkBackground)
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(Color.gray.opacity(0.2), lineWidth: 1)
                    )
            )
        }
    }
}

#Preview {
    ProfileView()
}
