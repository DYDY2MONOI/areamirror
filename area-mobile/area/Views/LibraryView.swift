//
//  LibraryView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct LibraryView: View {
    @StateObject private var authService = AuthService.shared
    @StateObject private var areaService = AreaService.shared
    @State private var selectedFilter = "All"
    @State private var selectedArea: Area?
    
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
                            
                            Text("\(areaService.userAreas.count) automation areas")
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
                        
                        if areaService.userAreas.isEmpty {
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
                                    ForEach($areaService.userAreas) { $area in
                                        AreaCard(
                                            area: $area,
                                            onEdit: {
                                                selectedArea = area
                                            },
                                            onDelete: {
                                            },
                                            onToggle: {
                                                try await areaService.toggleArea(areaId: area.id)
                                            }
                                        )
                                    }
                                }
                                .padding(.horizontal, 20)
                            }
                        }
                        
                        Spacer()
                    }
                }
            }
            .sheet(item: $selectedArea, onDismiss: {
                Task {
                    await areaService.fetchUserAreas()
                }
            }) { area in
                EditAreaView(area: area)
            }
        }
        .navigationTitle("Library")
        .navigationBarHidden(true)
        .onAppear {
            Task {
                await areaService.fetchUserAreas()
            }
        }
    }
    
    private func convertAreaToAreaItem(_ area: Area) -> AreaItem {
        return AreaItem(
            id: area.id,
            name: area.name,
            description: area.description,
            isActive: area.isActive,
            action: area.triggerService,
            reaction: area.actionService
        )
    }
    
    private func filterAreas() {
        print("Filtering by: \(selectedFilter)")
    }
}

struct AreaItem: Identifiable {
    let id: String
    let name: String
    let description: String
    let isActive: Bool
    let action: String
    let reaction: String
}

struct AreaCard: View {
    @Binding var area: Area
    var onEdit: () -> Void
    var onDelete: () -> Void
    var onToggle: () async throws -> Area
    @State private var isProcessing = false
    
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
                    guard !isProcessing else { return }
                    isProcessing = true
                    let originalState = area.isActive
                    Task { @MainActor in
                        defer { isProcessing = false }
                        do {
                            let updatedArea = try await onToggle()
                            area = updatedArea
                        } catch {
                            area.isActive = originalState
                            print("Error toggling area: \(error.localizedDescription)")
                        }
                    }
                }) {
                    if isProcessing {
                        ProgressView()
                            .progressViewStyle(CircularProgressViewStyle())
                            .frame(width: 24, height: 24)
                            .tint(area.isActive ? .green : .red)
                    } else {
                        Image(systemName: area.isActive ? "power" : "poweroff")
                            .font(.title2)
                            .foregroundColor(area.isActive ? .green : .red)
                    }
                }
                .buttonStyle(.plain)
                .disabled(isProcessing)
            }
            
            HStack {
                VStack(alignment: .leading, spacing: 4) {
                    Text("When: \(area.triggerService)")
                        .font(.caption)
                        .foregroundColor(.gray)
                    
                    Text("Then: \(area.actionService)")
                        .font(.caption)
                        .foregroundColor(.gray)
                }
                
                Spacer()
                
                HStack(spacing: 12) {
                    Button(action: onEdit) {
                        Image(systemName: "pencil")
                            .font(.title3)
                            .foregroundColor(AppColors.primaryBlue)
                    }
                    
                    Button(action: onDelete) {
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
