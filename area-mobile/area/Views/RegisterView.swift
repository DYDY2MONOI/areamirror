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
    @State private var showPassword = false
    @State private var showConfirmPassword = false
    @State private var acceptTerms = false
    
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
                            showPassword: $showPassword,
                            showConfirmPassword: $showConfirmPassword
                        )
                        
                        TermsAndConditions(acceptTerms: $acceptTerms)
                        
                        RegisterButton {
                            print("Register with email: \(email)")
                            onRegisterSuccess()
                        }
                        
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
    @Binding var showPassword: Bool
    @Binding var showConfirmPassword: Bool
    
    var body: some View {
        VStack(spacing: 20) {
            HStack(spacing: 12) {
                Image(systemName: "envelope")
                    .foregroundColor(.gray)
                    .frame(width: 20)
                
                TextField("Email", text: $email)
                    .textFieldStyle(PlainTextFieldStyle())
                    .foregroundColor(.white)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
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
                } else {
                    SecureField("Password", text: $password)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
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
                } else {
                    SecureField("Confirm Password", text: $confirmPassword)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
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
    
    var body: some View {
        Button(action: action) {
            Text("CREATE ACCOUNT")
                .font(AppTextStyles.button)
                .foregroundColor(.white)
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
