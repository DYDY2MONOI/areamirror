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
    @State private var showNewArea = false
    
    var body: some View {
        TabView(selection: $selectedTab) {
            HomeView(onLogout: {
            })
                .tabItem {
                    Image(systemName: selectedTab == 0 ? "house.fill" : "house")
                    Text("Home")
                }
                .tag(0)
            
            SearchView()
                .tabItem {
                    Image(systemName: selectedTab == 1 ? "magnifyingglass.circle.fill" : "magnifyingglass.circle")
                    Text("Search")
                }
                .tag(1)
            
            CreateView()
                .tabItem {
                    Image(systemName: selectedTab == 2 ? "plus.circle.fill" : "plus.circle")
                    Text("Create")
                }
                .tag(2)
            
            LibraryView()
                .tabItem {
                    Image(systemName: selectedTab == 3 ? "books.vertical.fill" : "books.vertical")
                    Text("Library")
                }
                .tag(3)
            
            ProfileView()
                .tabItem {
                    Image(systemName: selectedTab == 4 ? "person.circle.fill" : "person.circle")
                    Text("Profile")
                }
                .tag(4)
        }
        .accentColor(AppColors.primaryBlue)
        .onAppear {
            if authService.isAuthenticated {
                selectedTab = 0
            } else {
                selectedTab = 4
            }
        }
        .onChange(of: selectedTab) { _, newValue in
            if newValue == 2 {
                showNewArea = true
                selectedTab = 0
            }
        }
        .fullScreenCover(isPresented: $showNewArea) {
            NewAreaView()
        }
    }
}

#Preview {
    MainTabView()
}
