//
//  LoginView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct LoginView: View {
    @State private var email = ""
    @State private var password = ""
    @State private var rememberMe = false
    @State private var showPassword = false
    @State private var showTestView = false
    @State private var showAlert = false
    
    @StateObject private var authService = AuthService.shared
    private let oauthProviders = OAuthProvider.availableProviders
    
    let onLoginSuccess: () -> Void
    let onSignUpTap: () -> Void
    
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
                        
                        LoginHeader()
                        
                        LoginForm(
                            email: $email,
                            password: $password,
                            showPassword: $showPassword
                        )
                        
                        LoginOptions(rememberMe: $rememberMe)
                        
                        LoginButton {
                            performLogin()
                        }
                        .disabled(authService.isLoading)
                        
                        DividerWithText(text: "Or")
                        
                        OAuthLoginButtons(
                            providers: oauthProviders,
                            onProviderSelected: loginWithProvider
                        )
                        
                        SignUpPrompt {
                            onSignUpTap()
                        }
                    }
                }
            }
        }
        .fullScreenCover(isPresented: $showTestView) {
            TestView()
        }
        .alert("Error", isPresented: $showAlert) {
            Button("OK") { }
        } message: {
            Text(authService.errorMessage ?? "An error occurred")
        }
        .onChange(of: authService.isAuthenticated) { _, isAuthenticated in
            if isAuthenticated {
                onLoginSuccess()
            }
        }
        .onTapGesture {
            UIApplication.shared.sendAction(#selector(UIResponder.resignFirstResponder), to: nil, from: nil, for: nil)
        }
    }
    
    private func performLogin() {
        guard !email.isEmpty && !password.isEmpty else {
            authService.errorMessage = "Please fill in all fields"
            showAlert = true
            return
        }
        
        authService.login(email: email, password: password)
        
        if authService.errorMessage != nil {
            showAlert = true
        }
    }

    private func loginWithProvider(_ provider: OAuthProvider) {
        authService.login(with: provider)
    }
}

#Preview {
    LoginView(
        onLoginSuccess: {},
        onSignUpTap: {}
    )
}
