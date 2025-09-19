//
//  LibraryView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct LibraryView: View {
    @StateObject private var authService = AuthService.shared
    @State private var selectedFilter = "All"
    @State private var areas: [AreaItem] = []
    
    let filters = ["All", "Active", "Inactive", "Recent", "Favorites"]
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                if !authService.isAuthenticated {
                    VStack(spacing: 30) {
                        Image(systemName: "books.vertical")
                            .font(.system(size: 80))
                            .foregroundColor(.gray)
                        
                        Text("Login Required")
                            .font(.title)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        
                        Text("Please log in to view your automation areas")
                            .font(.body)
                            .foregroundColor(.gray)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal, 40)
                        
                        Button(action: {
                        }) {
                            Text("LOGIN")
                                .font(.headline)
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(AppGradients.button)
                                .cornerRadius(12)
                        }
                        .padding(.horizontal, 40)
                    }
                } else {
                    VStack(spacing: 20) {
                        VStack(spacing: 12) {
                            Text("My Areas")
                                .font(.title)
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                            
                            Text("\(areas.count) automation areas")
                                .font(.body)
                                .foregroundColor(.gray)
                        }
                        .padding(.top, 20)
                        
                        ScrollView(.horizontal, showsIndicators: false) {
                            HStack(spacing: 12) {
                                ForEach(filters, id: \.self) { filter in
                                    Button(action: {
                                        selectedFilter = filter
                                        filterAreas()
                                    }) {
                                        Text(filter)
                                            .font(.system(size: 14, weight: .medium))
                                            .foregroundColor(selectedFilter == filter ? .white : .gray)
                                            .padding(.horizontal, 16)
                                            .padding(.vertical, 8)
                                            .background(
                                                RoundedRectangle(cornerRadius: 20)
                                                    .fill(selectedFilter == filter ? AppColors.primaryBlue : Color.clear)
                                                    .overlay(
                                                        RoundedRectangle(cornerRadius: 20)
                                                            .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                                                    )
                                            )
                                    }
                                }
                            }
                            .padding(.horizontal, 20)
                        }
                        
                        if areas.isEmpty {
                            VStack(spacing: 20) {
                                Image(systemName: "plus.circle")
                                    .font(.system(size: 60))
                                    .foregroundColor(.gray)
                                
                                Text("No Areas Yet")
                                    .font(.title2)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                
                                Text("Create your first automation area to get started")
                                    .font(.body)
                                    .foregroundColor(.gray)
                                    .multilineTextAlignment(.center)
                                    .padding(.horizontal, 40)
                                
                                Button(action: {
                                }) {
                                    Text("CREATE AREA")
                                        .font(.headline)
                                        .foregroundColor(.white)
                                        .frame(maxWidth: .infinity)
                                        .padding(.vertical, 16)
                                        .background(AppGradients.button)
                                        .cornerRadius(12)
                                }
                                .padding(.horizontal, 40)
                            }
                            .padding(.top, 60)
                        } else {
                            ScrollView {
                                LazyVStack(spacing: 16) {
                                    ForEach(areas) { area in
                                        AreaCard(area: area)
                                    }
                                }
                                .padding(.horizontal, 20)
                            }
                        }
                        
                        Spacer()
                    }
                }
            }
        }
        .navigationTitle("Library")
        .navigationBarHidden(true)
        .onAppear {
            loadAreas()
        }
    }
    
    private func loadAreas() {
        areas = [
            AreaItem(id: 1, name: "Email to Slack", description: "Send Slack message when new email arrives", isActive: true, action: "New Email", reaction: "Send Slack Message"),
            AreaItem(id: 2, name: "GitHub to Discord", description: "Notify Discord when new commit", isActive: true, action: "New Commit", reaction: "Send Discord Message"),
            AreaItem(id: 3, name: "Weather Alert", description: "Send email when temperature is high", isActive: false, action: "High Temperature", reaction: "Send Email")
        ]
    }
    
    private func filterAreas() {
        print("Filtering by: \(selectedFilter)")
    }
}

struct AreaItem: Identifiable {
    let id: Int
    let name: String
    let description: String
    let isActive: Bool
    let action: String
    let reaction: String
}

struct AreaCard: View {
    let area: AreaItem
    @State private var isToggled = false
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text(area.name)
                        .font(.headline)
                        .foregroundColor(.white)
                    
                    Text(area.description)
                        .font(.caption)
                        .foregroundColor(.gray)
                }
                
                Spacer()
                
                Button(action: {
                    isToggled.toggle()
                }) {
                    Image(systemName: isToggled ? "power" : "poweroff")
                        .font(.title2)
                        .foregroundColor(isToggled ? .green : .red)
                }
            }
            
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("When: \(area.action)")
                        .font(.caption)
                        .foregroundColor(.gray)
                    
                    Text("Then: \(area.reaction)")
                        .font(.caption)
                        .foregroundColor(.gray)
                }
                
                Spacer()
                
                HStack(spacing: 12) {
                    Button(action: {
                    }) {
                        Image(systemName: "pencil")
                            .font(.title3)
                            .foregroundColor(AppColors.primaryBlue)
                    }
                    
                    Button(action: {
                    }) {
                        Image(systemName: "trash")
                            .font(.title3)
                            .foregroundColor(.red)
                    }
                }
            }
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(AppColors.darkBackground)
                .overlay(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(Color.gray.opacity(0.2), lineWidth: 1)
                )
        )
    }
}

#Preview {
    LibraryView()
}
