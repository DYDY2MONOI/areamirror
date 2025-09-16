//
//  ContentView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct ContentView: View {
    @State private var email = ""
    @State private var password = ""
    @State private var rememberMe = false
    @State private var showPassword = false
    @State private var animationOffset: CGFloat = 0
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                LinearGradient(
                    gradient: Gradient(colors: [
                        Color.black,
                        Color(red: 0.05, green: 0.05, blue: 0.15),
                        Color(red: 0.1, green: 0.1, blue: 0.2)
                    ]),
                    startPoint: .topLeading,
                    endPoint: .bottomTrailing
                )
                .ignoresSafeArea()
                
                VStack {
                    HStack {
                        Spacer()
                        RoundedRectangle(cornerRadius: 50)
                            .fill(
                                LinearGradient(
                                    gradient: Gradient(colors: [
                                        Color(red: 0.34, green: 0.50, blue: 0.91).opacity(0.3),
                                        Color(red: 0.52, green: 0.81, blue: 0.92).opacity(0.2)
                                    ]),
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 120, height: 80)
                            .rotationEffect(.degrees(15))
                            .offset(x: 20, y: -50)
                    }
                    Spacer()
                    HStack {
                        RoundedRectangle(cornerRadius: 40)
                            .fill(
                                LinearGradient(
                                    gradient: Gradient(colors: [
                                        Color(red: 0.53, green: 0.38, blue: 0.82).opacity(0.3),
                                        Color(red: 0.76, green: 0.78, blue: 0.89).opacity(0.2)
                                    ]),
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 100, height: 60)
                            .rotationEffect(.degrees(-20))
                            .offset(x: -30, y: 50)
                        Spacer()
                    }
                }
                
                ScrollView {
                    VStack(spacing: 0) {
                        Spacer(minLength: 60)
                        
                        Text("AREA")
                            .font(.system(size: 36, weight: .bold, design: .rounded))
                            .foregroundColor(.white)
                            .padding(.bottom, 60)
                        
                        VStack(spacing: 12) {
                            Text("LOGIN TO YOUR ACCOUNT")
                                .font(.system(size: 24, weight: .bold))
                                .foregroundColor(.white)
                                .multilineTextAlignment(.center)
                            
                            Text("Enter your login information")
                                .font(.system(size: 16))
                                .foregroundColor(.gray)
                                .multilineTextAlignment(.center)
                        }
                        .padding(.bottom, 40)
                        
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
                                    .fill(Color(red: 0.1, green: 0.1, blue: 0.15))
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
                                    .fill(Color(red: 0.1, green: 0.1, blue: 0.15))
                                    .overlay(
                                        RoundedRectangle(cornerRadius: 12)
                                            .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                                    )
                            )
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 20)
                        
                        HStack {
                            Button(action: {
                                rememberMe.toggle()
                            }) {
                                HStack(spacing: 8) {
                                    Image(systemName: rememberMe ? "checkmark.square.fill" : "square")
                                        .foregroundColor(rememberMe ? Color(red: 0.34, green: 0.50, blue: 0.91) : .gray)
                                    
                                    Text("Remember me")
                                        .foregroundColor(.gray)
                                        .font(.system(size: 14))
                                }
                            }
                            
                            Spacer()
                            
                            Button("Forgot password?") {
                            }
                            .foregroundColor(Color(red: 0.34, green: 0.50, blue: 0.91))
                            .font(.system(size: 14))
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 30)
                        
                        Button(action: {
                            print("Login with email: \(email)")
                        }) {
                            Text("LOGIN")
                                .font(.system(size: 16, weight: .bold))
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(
                                    LinearGradient(
                                        gradient: Gradient(colors: [
                                            Color(red: 0.34, green: 0.50, blue: 0.91),
                                            Color(red: 0.53, green: 0.38, blue: 0.82)
                                        ]),
                                        startPoint: .leading,
                                        endPoint: .trailing
                                    )
                                )
                                .cornerRadius(12)
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 30)
                        
                        HStack {
                            Rectangle()
                                .fill(Color.gray.opacity(0.3))
                                .frame(height: 1)
                            
                            Text("Or")
                                .foregroundColor(.gray)
                                .font(.system(size: 14))
                                .padding(.horizontal, 16)
                            
                            Rectangle()
                                .fill(Color.gray.opacity(0.3))
                                .frame(height: 1)
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 30)
                        
                        VStack(spacing: 12) {
                            Button(action: {
                                print("Google Login")
                            }) {
                                HStack(spacing: 12) {
                                    Image(systemName: "globe")
                                        .font(.system(size: 18))
                                        .foregroundColor(.red)
                                    
                                    Text("GOOGLE")
                                        .font(.system(size: 16, weight: .semibold))
                                        .foregroundColor(.black)
                                }
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(Color.white)
                                .cornerRadius(12)
                            }
                            
                            Button(action: {
                                print("Apple Login")
                            }) {
                                HStack(spacing: 12) {
                                    Image(systemName: "applelogo")
                                        .font(.system(size: 18))
                                        .foregroundColor(.black)
                                    
                                    Text("APPLE")
                                        .font(.system(size: 16, weight: .semibold))
                                        .foregroundColor(.black)
                                }
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(Color.white)
                                .cornerRadius(12)
                            }
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 40)
                        
                        HStack {
                            Text("Don't have an account? ")
                                .foregroundColor(.gray)
                                .font(.system(size: 14))
                            
                            Button("Sign up") {
                                print("Sign up")
                            }
                            .foregroundColor(Color(red: 0.34, green: 0.50, blue: 0.91))
                            .font(.system(size: 14, weight: .semibold))
                        }
                        .padding(.bottom, 40)
                    }
                }
            }
        }
    }
}

#Preview {
    ContentView()
}
