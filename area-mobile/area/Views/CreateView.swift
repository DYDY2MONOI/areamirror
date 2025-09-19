//
//  CreateView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct CreateView: View {
    @StateObject private var authService = AuthService.shared
    @State private var selectedAction: String = ""
    @State private var selectedReaction: String = ""
    @State private var areaName: String = ""
    @State private var areaDescription: String = ""
    @State private var showActionPicker = false
    @State private var showReactionPicker = false
    
    let actions = ["New Email", "New Commit", "High Temperature", "New Message", "File Upload"]
    let reactions = ["Send Slack Message", "Send Email", "Create Issue", "Post Tweet", "Save File"]
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                if !authService.isAuthenticated {
                    VStack(spacing: 30) {
                        Image(systemName: "lock.circle")
                            .font(.system(size: 80))
                            .foregroundColor(.gray)
                        
                        Text("Login Required")
                            .font(.title)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                        
                        Text("Please log in to create automation areas")
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
                    ScrollView {
                        VStack(spacing: 24) {
                            VStack(spacing: 12) {
                                Text("Create New Area")
                                    .font(.title)
                                    .fontWeight(.bold)
                                    .foregroundColor(.white)
                                
                                Text("Build your automation workflow")
                                    .font(.body)
                                    .foregroundColor(.gray)
                            }
                            .padding(.top, 20)
                            
                            VStack(alignment: .leading, spacing: 12) {
                                Text("When this happens...")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                
                                Button(action: {
                                    showActionPicker = true
                                }) {
                                    HStack {
                                        Text(selectedAction.isEmpty ? "Select Action" : selectedAction)
                                            .foregroundColor(selectedAction.isEmpty ? .gray : .white)
                                        
                                        Spacer()
                                        
                                        Image(systemName: "chevron.down")
                                            .foregroundColor(.gray)
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
                            }
                            
                            Image(systemName: "arrow.down")
                                .font(.title2)
                                .foregroundColor(AppColors.primaryBlue)
                                .padding(.vertical, 8)
                            
                            VStack(alignment: .leading, spacing: 12) {
                                Text("Then do this...")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                
                                Button(action: {
                                    showReactionPicker = true
                                }) {
                                    HStack {
                                        Text(selectedReaction.isEmpty ? "Select Reaction" : selectedReaction)
                                            .foregroundColor(selectedReaction.isEmpty ? .gray : .white)
                                        
                                        Spacer()
                                        
                                        Image(systemName: "chevron.down")
                                            .foregroundColor(.gray)
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
                            }
                            
                            VStack(alignment: .leading, spacing: 12) {
                                Text("Area Name")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                
                                TextField("My Automation Area", text: $areaName)
                                    .textFieldStyle(PlainTextFieldStyle())
                                    .foregroundColor(.white)
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
                            
                            VStack(alignment: .leading, spacing: 12) {
                                Text("Description (Optional)")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                
                                TextField("Describe what this area does...", text: $areaDescription, axis: .vertical)
                                    .textFieldStyle(PlainTextFieldStyle())
                                    .foregroundColor(.white)
                                    .lineLimit(3...6)
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
                            
                            Button(action: {
                                createArea()
                            }) {
                                Text("CREATE AREA")
                                    .font(.headline)
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .padding(.vertical, 16)
                                    .background(
                                        selectedAction.isEmpty || selectedReaction.isEmpty || areaName.isEmpty ?
                                        LinearGradient(gradient: Gradient(colors: [Color.gray, Color.gray]), startPoint: .leading, endPoint: .trailing) : AppGradients.button
                                    )
                                    .cornerRadius(12)
                            }
                            .disabled(selectedAction.isEmpty || selectedReaction.isEmpty || areaName.isEmpty)
                            .padding(.top, 20)
                        }
                        .padding(.horizontal, 20)
                        .padding(.bottom, 40)
                    }
                }
            }
        }
        .navigationTitle("Create")
        .navigationBarHidden(true)
        .sheet(isPresented: $showActionPicker) {
            ActionPickerView(selectedAction: $selectedAction, actions: actions)
        }
        .sheet(isPresented: $showReactionPicker) {
            ReactionPickerView(selectedReaction: $selectedReaction, reactions: reactions)
        }
    }
    
    private func createArea() {
        print("Creating area: \(areaName)")
        print("Action: \(selectedAction)")
        print("Reaction: \(selectedReaction)")
    }
}

struct ActionPickerView: View {
    @Binding var selectedAction: String
    let actions: [String]
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                List(actions, id: \.self) { action in
                    Button(action: {
                        selectedAction = action
                        dismiss()
                    }) {
                        HStack {
                            Text(action)
                                .foregroundColor(.white)
                            
                            Spacer()
                            
                            if selectedAction == action {
                                Image(systemName: "checkmark")
                                    .foregroundColor(AppColors.primaryBlue)
                            }
                        }
                        .padding(.vertical, 8)
                    }
                    .listRowBackground(Color.clear)
                }
                .listStyle(PlainListStyle())
            }
            .navigationTitle("Select Action")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Cancel") {
                        dismiss()
                    }
                    .foregroundColor(AppColors.primaryBlue)
                }
            }
        }
    }
}

struct ReactionPickerView: View {
    @Binding var selectedReaction: String
    let reactions: [String]
    @Environment(\.dismiss) private var dismiss
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                List(reactions, id: \.self) { reaction in
                    Button(action: {
                        selectedReaction = reaction
                        dismiss()
                    }) {
                        HStack {
                            Text(reaction)
                                .foregroundColor(.white)
                            
                            Spacer()
                            
                            if selectedReaction == reaction {
                                Image(systemName: "checkmark")
                                    .foregroundColor(AppColors.primaryBlue)
                            }
                        }
                        .padding(.vertical, 8)
                    }
                    .listRowBackground(Color.clear)
                }
                .listStyle(PlainListStyle())
            }
            .navigationTitle("Select Reaction")
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Cancel") {
                        dismiss()
                    }
                    .foregroundColor(AppColors.primaryBlue)
                }
            }
        }
    }
}

#Preview {
    CreateView()
}
