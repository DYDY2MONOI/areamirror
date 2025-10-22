//
//  RegisterView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct RegisterView: View {
    @State private var email = ""
    @State private var password = ""
    @State private var confirmPassword = ""
    @State private var firstName = ""
    @State private var lastName = ""
    @State private var showPassword = false
    @State private var showConfirmPassword = false
    @State private var acceptTerms = false
    @State private var showAlert = false
    
    @StateObject private var authService = AuthService.shared
    
    let onRegisterSuccess: () -> Void
    let onLoginTap: () -> Void
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                DecorativeShapes()
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer(minLength: 60)
                        
                        Text("AREA")
                            .font(AppTextStyles.title)
                            .foregroundColor(.white)
                            .padding(.bottom, 60)
                        
                        RegisterHeader()
                        
                        RegisterForm(
                            email: $email,
                            password: $password,
                            confirmPassword: $confirmPassword,
                            firstName: $firstName,
                            lastName: $lastName,
                            showPassword: $showPassword,
                            showConfirmPassword: $showConfirmPassword
                        )
                        
                        TermsAndConditions(acceptTerms: $acceptTerms)
                        
                        RegisterButton {
                            performRegister()
                        }
                        .disabled(authService.isLoading)
                        
                        DividerWithText(text: "Or")
                        
                        SocialLoginButtons(
                            onGoogleLogin: {
                                print("Google Register")
                            },
                            onAppleLogin: {
                                print("Apple Register")
                            }
                        )
                        
                        LoginPrompt {
                            onLoginTap()
                        }
                    }
                }
            }
        }
        .alert("Error", isPresented: $showAlert) {
            Button("OK") { }
        } message: {
            Text(authService.errorMessage ?? "An error occurred")
        }
        .onChange(of: authService.isAuthenticated) { _, isAuthenticated in
            if isAuthenticated {
                onRegisterSuccess()
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
    }
    
    private func hideKeyboard() {
        UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
    }
    
    private func performRegister() {
        guard !email.isEmpty && !password.isEmpty else {
            authService.errorMessage = "Please fill in all required fields"
            showAlert = true
            return
        }
        
        guard password == confirmPassword else {
            authService.errorMessage = "Passwords do not match"
            showAlert = true
            return
        }
        
        guard acceptTerms else {
            authService.errorMessage = "Please accept the terms and conditions"
            showAlert = true
            return
        }
        
        authService.register(
            email: email,
            password: password,
            firstName: firstName.isEmpty ? nil : firstName,
            lastName: lastName.isEmpty ? nil : lastName
        )
        
        if authService.errorMessage != nil {
            showAlert = true
        }
    }
}

struct RegisterHeader: View {
    var body: some View {
        VStack(spacing: 12) {
            Text("CREATE YOUR ACCOUNT")
                .font(AppTextStyles.subtitle)
                .foregroundColor(.white)
                .multilineTextAlignment(.center)
            
            Text("Enter your information to get started")
                .font(AppTextStyles.body)
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)
        }
        .padding(.bottom, 40)
    }
}

struct RegisterForm: View {
    @Binding var email: String
    @Binding var password: String
    @Binding var confirmPassword: String
    @Binding var firstName: String
    @Binding var lastName: String
    @Binding var showPassword: Bool
    @Binding var showConfirmPassword: Bool
    
    var body: some View {
        VStack(spacing: 20) {
            HStack(spacing: 12) {
                HStack(spacing: 12) {
                    Image(systemName: "person")
                        .foregroundColor(.gray)
                        .frame(width: 20)
                    
                    TextField("First Name", text: $firstName)
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
                
                HStack(spacing: 12) {
                    Image(systemName: "person")
                        .foregroundColor(.gray)
                        .frame(width: 20)
                    
                    TextField("Last Name", text: $lastName)
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
            
            HStack(spacing: 12) {
                Image(systemName: "envelope")
                    .foregroundColor(.gray)
                    .frame(width: 20)
                
                TextField("Email", text: $email)
                    .textFieldStyle(PlainTextFieldStyle())
                    .foregroundColor(.white)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
                    .keyboardType(.emailAddress)
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
            
            HStack(spacing: 12) {
                Image(systemName: "lock")
                    .foregroundColor(.gray)
                    .frame(width: 20)
                
                if showPassword {
                    TextField("Password", text: $password)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                        .submitLabel(.next)
                } else {
                    SecureField("Password", text: $password)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                        .submitLabel(.next)
                }
                
                Button(action: {
                    showPassword.toggle()
                }) {
                    Image(systemName: showPassword ? "eye" : "eye.slash")
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
            
            HStack(spacing: 12) {
                Image(systemName: "lock")
                    .foregroundColor(.gray)
                    .frame(width: 20)
                
                if showConfirmPassword {
                    TextField("Confirm Password", text: $confirmPassword)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                        .submitLabel(.done)
                } else {
                    SecureField("Confirm Password", text: $confirmPassword)
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
        .padding(.horizontal, 24)
        .padding(.bottom, 20)
    }
}

struct TermsAndConditions: View {
    @Binding var acceptTerms: Bool
    
    var body: some View {
        HStack(alignment: .top, spacing: 8) {
            Button(action: {
                acceptTerms.toggle()
            }) {
                Image(systemName: acceptTerms ? "checkmark.square.fill" : "square")
                    .foregroundColor(acceptTerms ? AppColors.primaryBlue : .gray)
            }
            
            Text("I agree to the Terms and Conditions and Privacy Policy")
                .foregroundColor(.gray)
                .font(AppTextStyles.caption)
                .multilineTextAlignment(.leading)
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 30)
    }
}

struct RegisterButton: View {
    let action: () -> Void
    @StateObject private var authService = AuthService.shared
    
    var body: some View {
        Button(action: action) {
            HStack {
                if authService.isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        .scaleEffect(0.8)
                }
                Text(authService.isLoading ? "CREATING..." : "CREATE ACCOUNT")
                    .font(AppTextStyles.button)
                    .foregroundColor(.white)
            }
            .frame(maxWidth: .infinity)
            .padding(.vertical, 16)
            .background(AppGradients.button)
            .cornerRadius(12)
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 30)
    }
}

struct LoginPrompt: View {
    let onLogin: () -> Void
    
    var body: some View {
        HStack {
            Text("Already have an account? ")
                .foregroundColor(.gray)
                .font(AppTextStyles.caption)
            
            Button("Login", action: onLogin)
                .foregroundColor(AppColors.primaryBlue)
                .font(.system(size: 14, weight: .semibold))
        }
        .padding(.bottom, 40)
    }
}

#Preview {
    RegisterView(
        onRegisterSuccess: {},
        onLoginTap: {}
    )
}
