//
//  EditProfileView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct EditProfileView: View {
    @StateObject private var authService = AuthService.shared
    @State private var firstName = ""
    @State private var lastName = ""
    @State private var showAlert = false
    @State private var showSuccessAlert = false
    
    let onDismiss: () -> Void
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        VStack(spacing: 16) {
                            Circle()
                                .fill(AppGradients.button)
                                .frame(width: 80, height: 80)
                                .overlay(
                                    Text(authService.currentUser?.firstName?.prefix(1).uppercased() ?? "U")
                                        .font(.title)
                                        .fontWeight(.bold)
                                        .foregroundColor(.white)
                                )
                            
                            Text("Edit Profile")
                                .font(AppTextStyles.title)
                                .foregroundColor(.white)
                                .multilineTextAlignment(.center)
                        }
                        .padding(.top, 20)
                        .padding(.bottom, 40)
                        
                        VStack(spacing: 20) {
                            VStack(alignment: .leading, spacing: 8) {
                                Text("First Name")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "person")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    TextField("Enter your first name", text: $firstName)
                                        .textFieldStyle(PlainTextFieldStyle())
                                        .foregroundColor(.white)
                                        .autocapitalization(.words)
                                        .submitLabel(.next)
                                }
                                .padding(.horizontal, 16)
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
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Last Name")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "person")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    TextField("Enter your last name", text: $lastName)
                                        .textFieldStyle(PlainTextFieldStyle())
                                        .foregroundColor(.white)
                                        .autocapitalization(.words)
                                        .submitLabel(.done)
                                }
                                .padding(.horizontal, 16)
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
                            
                            VStack(alignment: .leading, spacing: 8) {
                                Text("Email")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "envelope")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    Text(authService.currentUser?.email ?? "")
                                        .foregroundColor(.gray)
                                        .padding(.horizontal, 16)
                                        .padding(.vertical, 16)
                                }
                                .frame(maxWidth: .infinity, alignment: .leading)
                                .background(
                                    RoundedRectangle(cornerRadius: 12)
                                        .fill(AppColors.darkBackground.opacity(0.5))
                                        .overlay(
                                            RoundedRectangle(cornerRadius: 12)
                                                .stroke(Color.gray.opacity(0.2), lineWidth: 1)
                                        )
                                )
                            }
                            
                            HStack {
                                Image(systemName: "info.circle")
                                    .foregroundColor(.blue)
                                Text("Email cannot be changed")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                            }
                            .padding(.top, 8)
                        }
                        .padding(.horizontal, 24)
                        
                        Button(action: {
                            saveProfile()
                        }) {
                            HStack {
                                if authService.isLoading {
                                    ProgressView()
                                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                        .scaleEffect(0.8)
                                }
                                Text(authService.isLoading ? "SAVING..." : "SAVE CHANGES")
                                    .font(AppTextStyles.button)
                                    .foregroundColor(.white)
                            }
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(AppGradients.button)
                            .cornerRadius(12)
                        }
                        .disabled(authService.isLoading)
                        .padding(.horizontal, 24)
                        .padding(.top, 40)
                        .padding(.bottom, 20)
                    }
                }
            }
        }
        .navigationTitle("Edit Profile")
        .navigationBarTitleDisplayMode(.inline)
        .navigationBarBackButtonHidden(true)
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                Button("Cancel") {
                    onDismiss()
                }
                .foregroundColor(AppColors.primaryBlue)
            }
        }
        .onAppear {
            loadCurrentProfile()
        }
        .onChange(of: authService.errorMessage) { errorMessage in
            if let error = errorMessage {
                showAlert = true
            }
        }
        .onChange(of: authService.isLoading) { isLoading in
            if !isLoading && authService.errorMessage == nil {
                showSuccessAlert = true
            }
        }
        .onTapGesture {
            hideKeyboard()
        }
        .gesture(
            DragGesture(minimumDistance: 0)
                .onEnded { _ in
                    hideKeyboard()
                }
        )
        .alert("Error", isPresented: $showAlert) {
            Button("OK") { }
        } message: {
            Text(authService.errorMessage ?? "An error occurred")
        }
        .alert("Success", isPresented: $showSuccessAlert) {
            Button("OK") {
                onDismiss()
            }
        } message: {
            Text("Profile updated successfully!")
        }
    }
    
    private func loadCurrentProfile() {
        firstName = authService.currentUser?.firstName ?? ""
        lastName = authService.currentUser?.lastName ?? ""
    }
    
    private func saveProfile() {
        guard !firstName.isEmpty || !lastName.isEmpty else {
            authService.errorMessage = "Please enter at least one name"
            showAlert = true
            return
        }
        
        authService.updateProfile(
            firstName: firstName.isEmpty ? nil : firstName,
            lastName: lastName.isEmpty ? nil : lastName
        )
    }
    
    private func hideKeyboard() {
        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
    }
}

#Preview {
    EditProfileView(onDismiss: {})
}
