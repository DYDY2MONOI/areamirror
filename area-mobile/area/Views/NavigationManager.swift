//
//  NavigationManager.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

enum AppScreen {
    case login
    case register
    case home
}

class NavigationManager: ObservableObject {
    @Published var currentScreen: AppScreen = .login
    @Published var isLoggedIn: Bool = false
    
    func navigateToLogin() {
        currentScreen = .login
        isLoggedIn = false
    }
    
    func navigateToRegister() {
        currentScreen = .register
        isLoggedIn = false
    }
    
    func navigateToHome() {
        currentScreen = .home
        isLoggedIn = true
    }
    
    func logout() {
        navigateToLogin()
    }
}

struct AppNavigationView: View {
    @StateObject private var navigationManager = NavigationManager()
    
    var body: some View {
        Group {
            switch navigationManager.currentScreen {
            case .login:
                LoginView(
                    onLoginSuccess: {
                        navigationManager.navigateToHome()
                    },
                    onSignUpTap: {
                        navigationManager.navigateToRegister()
                    }
                )
                
            case .register:
                RegisterView(
                    onRegisterSuccess: {
                        navigationManager.navigateToHome()
                    },
                    onLoginTap: {
                        navigationManager.navigateToLogin()
                    }
                )
                
            case .home:
                HomeView(
                    onLogout: {
                        navigationManager.logout()
                    }
                )
            }
        }
        .environmentObject(navigationManager)
    }
}

#Preview {
    AppNavigationView()
}
