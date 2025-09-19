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
                            print("Login with email: \(email)")
                            onLoginSuccess()
                        }
                        
                        DividerWithText(text: "Or")
                        
                        SocialLoginButtons(
                            onGoogleLogin: {
                                print("Google Login")
                            },
                            onAppleLogin: {
                                print("Apple Login")
                            }
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
    }
}

#Preview {
    LoginView(
        onLoginSuccess: {},
        onSignUpTap: {}
    )
}
