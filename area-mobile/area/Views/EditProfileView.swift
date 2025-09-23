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
    @State private var phone = ""
    @State private var country = ""
    @State private var currentPassword = ""
    @State private var newPassword = ""
    @State private var confirmPassword = ""
    @State private var showCurrentPassword = false
    @State private var showNewPassword = false
    @State private var showConfirmPassword = false
    @State private var showAlert = false
    @State private var showSuccessAlert = false
    @State private var isChangingPassword = false
    
    let onDismiss: () -> Void
    
    var body: some View {
        ZStack {
            AppGradients.background
                .ignoresSafeArea()
            
            VStack(spacing: 0) {
                HStack {
                    Button(action: {
                        onDismiss()
                    }) {
                        Image(systemName: "xmark")
                            .font(.title2)
                            .fontWeight(.medium)
                            .foregroundColor(.white)
                            .frame(width: 32, height: 32)
                            .background(
                                Circle()
                                    .fill(Color.black.opacity(0.3))
                            )
                    }
                    
                    Spacer()
                    
                    Text("Edit Profile")
                        .font(AppTextStyles.title)
                        .foregroundColor(.white)
                    
                    Spacer()
                    
                    Color.clear
                        .frame(width: 32, height: 32)
                }
                .padding(.horizontal, 20)
                .padding(.top, 10)
                .padding(.bottom, 20)
                
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
                                        .disableAutocorrection(true)
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
                                        .disableAutocorrection(true)
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
                                Text("Phone")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "phone")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    TextField("Enter your phone number", text: $phone)
                                        .textFieldStyle(PlainTextFieldStyle())
                                        .foregroundColor(.white)
                                        .keyboardType(.phonePad)
                                        .autocapitalization(.none)
                                        .disableAutocorrection(true)
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
                                Text("Country")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "globe")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    TextField("Enter your country", text: $country)
                                        .textFieldStyle(PlainTextFieldStyle())
                                        .foregroundColor(.white)
                                        .autocapitalization(.words)
                                        .disableAutocorrection(true)
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
                                Text("Email")
                                    .font(AppTextStyles.caption)
                                    .foregroundColor(.gray)
                                
                                HStack(spacing: 12) {
                                    Image(systemName: "envelope")
                                        .foregroundColor(.gray)
                                        .frame(width: 20)
                                    
                                    Text(authService.currentUser?.email ?? "")
                                        .foregroundColor(.gray)
                                }
                                .frame(maxWidth: .infinity, alignment: .leading)
                                .padding(.horizontal, 16)
                                .padding(.vertical, 16)
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
                            
                            Divider()
                                .background(Color.gray.opacity(0.3))
                                .padding(.vertical, 20)
                            
                            VStack(alignment: .leading, spacing: 16) {
                                HStack {
                                    Text("Change Password")
                                        .font(AppTextStyles.subtitle)
                                        .foregroundColor(.white)
                                    
                                    Spacer()
                                    
                                    Button(action: {
                                        isChangingPassword.toggle()
                                    }) {
                                        Text(isChangingPassword ? "Cancel" : "Change")
                                            .font(AppTextStyles.caption)
                                            .foregroundColor(AppColors.primaryBlue)
                                    }
                                }
                                
                                if isChangingPassword {
                                    VStack(spacing: 16) {
                                        VStack(alignment: .leading, spacing: 8) {
                                            Text("Current Password")
                                                .font(AppTextStyles.caption)
                                                .foregroundColor(.gray)
                                            
                                            HStack(spacing: 12) {
                                                Image(systemName: "lock")
                                                    .foregroundColor(.gray)
                                                    .frame(width: 20)
                                                
                                                if showCurrentPassword {
                                                    TextField("Current password", text: $currentPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.next)
                                                } else {
                                                    SecureField("Current password", text: $currentPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.next)
                                                }
                                                
                                                Button(action: {
                                                    showCurrentPassword.toggle()
                                                }) {
                                                    Image(systemName: showCurrentPassword ? "eye" : "eye.slash")
                                                        .foregroundColor(.gray)
                                                }
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
                                            Text("New Password")
                                                .font(AppTextStyles.caption)
                                                .foregroundColor(.gray)
                                            
                                            HStack(spacing: 12) {
                                                Image(systemName: "lock")
                                                    .foregroundColor(.gray)
                                                    .frame(width: 20)
                                                
                                                if showNewPassword {
                                                    TextField("New password", text: $newPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.next)
                                                } else {
                                                    SecureField("New password", text: $newPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.next)
                                                }
                                                
                                                Button(action: {
                                                    showNewPassword.toggle()
                                                }) {
                                                    Image(systemName: showNewPassword ? "eye" : "eye.slash")
                                                        .foregroundColor(.gray)
                                                }
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
                                            Text("Confirm New Password")
                                                .font(AppTextStyles.caption)
                                                .foregroundColor(.gray)
                                            
                                            HStack(spacing: 12) {
                                                Image(systemName: "lock")
                                                    .foregroundColor(.gray)
                                                    .frame(width: 20)
                                                
                                                if showConfirmPassword {
                                                    TextField("Confirm new password", text: $confirmPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.done)
                                                } else {
                                                    SecureField("Confirm new password", text: $confirmPassword)
                                                        .textFieldStyle(PlainTextFieldStyle())
                                                        .foregroundColor(.white)
                                                        .autocapitalization(.none)
                                                        .disableAutocorrection(true)
                                                        .submitLabel(.done)
                                                }
                                                
                                                Button(action: {
                                                    showConfirmPassword.toggle()
                                                }) {
                                                    Image(systemName: showConfirmPassword ? "eye" : "eye.slash")
                                                        .foregroundColor(.gray)
                                                }
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
                                    }
                                }
                            }
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
        phone = authService.currentUser?.phone ?? ""
        country = authService.currentUser?.country ?? ""
    }
    
    private func saveProfile() {
        guard !firstName.isEmpty || !lastName.isEmpty else {
            authService.errorMessage = "Please enter at least one name"
            showAlert = true
            return
        }
        
        if isChangingPassword {
            guard !currentPassword.isEmpty && !newPassword.isEmpty && !confirmPassword.isEmpty else {
                authService.errorMessage = "Please fill in all password fields"
                showAlert = true
                return
            }
            
            guard newPassword == confirmPassword else {
                authService.errorMessage = "New passwords do not match"
                showAlert = true
                return
            }
            
            guard newPassword.count >= 6 else {
                authService.errorMessage = "New password must be at least 6 characters"
                showAlert = true
                return
            }
        }
        
        authService.updateProfile(
            firstName: firstName.isEmpty ? nil : firstName,
            lastName: lastName.isEmpty ? nil : lastName,
            phone: phone.isEmpty ? nil : phone,
            country: country.isEmpty ? nil : country,
            currentPassword: isChangingPassword ? currentPassword : nil,
            newPassword: isChangingPassword ? newPassword : nil
        )
        
        if isChangingPassword {
            currentPassword = ""
            newPassword = ""
            confirmPassword = ""
            isChangingPassword = false
        }
    }
    
    private func hideKeyboard() {
        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
    }
}

#Preview {
    EditProfileView(onDismiss: {})
}