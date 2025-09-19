//
//  MainTabView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct MainTabView: View {
    @StateObject private var authService = AuthService.shared
    @State private var selectedTab = 0
    
    var body: some View {
        TabView(selection: $selectedTab) {
            // Accueil
            HomeView(onLogout: {
                // Gérer la déconnexion si nécessaire
            })
                .tabItem {
                    Image(systemName: selectedTab == 0 ? "house.fill" : "house")
                    Text("Home")
                }
                .tag(0)
            
            // Recherche
            SearchView()
                .tabItem {
                    Image(systemName: selectedTab == 1 ? "magnifyingglass.circle.fill" : "magnifyingglass.circle")
                    Text("Search")
                }
                .tag(1)
            
            // Créer
            CreateView()
                .tabItem {
                    Image(systemName: selectedTab == 2 ? "plus.circle.fill" : "plus.circle")
                    Text("Create")
                }
                .tag(2)
            
            // Bibliothèque
            LibraryView()
                .tabItem {
                    Image(systemName: selectedTab == 3 ? "books.vertical.fill" : "books.vertical")
                    Text("Library")
                }
                .tag(3)
            
            // Profil/Connexion
            ProfileView()
                .tabItem {
                    Image(systemName: selectedTab == 4 ? "person.circle.fill" : "person.circle")
                    Text("Profile")
                }
                .tag(4)
        }
        .accentColor(AppColors.primaryBlue)
        .onAppear {
            // Vérifier l'état d'authentification au démarrage
            if authService.isAuthenticated {
                selectedTab = 0 // Aller à l'accueil si connecté
            } else {
                selectedTab = 4 // Aller au profil pour se connecter
            }
        }
    }
}

#Preview {
    MainTabView()
}
