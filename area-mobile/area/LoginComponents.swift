//
//  LoginComponents.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct LoginHeader: View {
    var body: some View {
        VStack(spacing: 12) {
            Text("LOGIN TO YOUR ACCOUNT")
                .font(AppTextStyles.subtitle)
                .foregroundColor(.white)
                .multilineTextAlignment(.center)
            
            Text("Enter your login information")
                .font(AppTextStyles.body)
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)
        }
        .padding(.bottom, 40)
    }
}

struct LoginForm: View {
    @Binding var email: String
    @Binding var password: String
    @Binding var showPassword: Bool
    
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
                        .submitLabel(.done)
                } else {
                    SecureField("Password", text: $password)
                        .textFieldStyle(PlainTextFieldStyle())
                        .foregroundColor(.white)
                        .autocapitalization(.none)
                        .disableAutocorrection(true)
                        .submitLabel(.done)
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
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 20)
    }
}

struct LoginOptions: View {
    @Binding var rememberMe: Bool
    
    var body: some View {
        HStack {
            Button(action: {
                rememberMe.toggle()
            }) {
                HStack(spacing: 8) {
                    Image(systemName: rememberMe ? "checkmark.square.fill" : "square")
                        .foregroundColor(rememberMe ? AppColors.primaryBlue : .gray)
                    
                    Text("Remember me")
                        .foregroundColor(.gray)
                        .font(AppTextStyles.caption)
                }
            }
            
            Spacer()
            
            Button("Forgot password?") {
            }
            .foregroundColor(AppColors.primaryBlue)
            .font(AppTextStyles.caption)
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 30)
    }
}

struct LoginButton: View {
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
                Text(authService.isLoading ? "LOGGING IN..." : "LOGIN")
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

struct OAuthLoginButtons: View {
    let providers: [OAuthProvider]
    let onProviderSelected: (OAuthProvider) -> Void

    var body: some View {
        VStack(spacing: 12) {
            ForEach(providers.filter { $0.isEnabled }) { provider in
                Button(action: {
                    onProviderSelected(provider)
                }) {
                    HStack(spacing: 12) {
                        Image(provider.iconName)
                            .resizable()
                            .renderingMode(.template)
                            .frame(width: 20, height: 20)
                            .foregroundColor(.white)

                        Text(provider.name.uppercased())
                            .font(.system(size: 16, weight: .semibold))
                            .foregroundColor(.white)
                    }
                    .frame(maxWidth: .infinity)
                    .padding(.vertical, 16)
                    .background(provider.color)
                    .cornerRadius(12)
                }
            }
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 40)
    }
}

struct SignUpPrompt: View {
    let onSignUp: () -> Void
    
    var body: some View {
        HStack {
            Text("Don't have an account? ")
                .foregroundColor(.gray)
                .font(AppTextStyles.caption)
            
            Button("Sign up", action: onSignUp)
                .foregroundColor(AppColors.primaryBlue)
                .font(.system(size: 14, weight: .semibold))
        }
        .padding(.bottom, 40)
    }
}

struct DividerWithText: View {
    let text: String
    
    var body: some View {
        HStack {
            Rectangle()
                .fill(Color.gray.opacity(0.3))
                .frame(height: 1)
            
            Text(text)
                .foregroundColor(.gray)
                .font(AppTextStyles.caption)
                .padding(.horizontal, 16)
            
            Rectangle()
                .fill(Color.gray.opacity(0.3))
                .frame(height: 1)
        }
        .padding(.horizontal, 24)
        .padding(.bottom, 30)
    }
}
