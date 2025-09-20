//
//  NewAreaView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct NewAreaView: View {
    @Environment(\.presentationMode) var presentationMode
    @State private var areaName = ""
    @State private var areaDescription = ""
    @State private var selectedActionService: Service?
    @State private var selectedReactionService: Service?
    @State private var showingActionSelection = false
    @State private var showingReactionSelection = false
    @State private var showingSuccessAlert = false
    
    var body: some View {
        NavigationView {
            ZStack {
                Color.black
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 30) {
                        VStack(spacing: 16) {
                            Text("Create New AREA")
                                .font(.system(size: 28, weight: .bold))
                                .foregroundColor(.white)
                            
                            Text("Connect two services to automate your workflow")
                                .font(.system(size: 16))
                                .foregroundColor(.gray)
                                .multilineTextAlignment(.center)
                                .padding(.horizontal, 20)
                        }
                        .padding(.top, 20)
                        
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Area Name")
                                .font(.system(size: 18, weight: .semibold))
                                .foregroundColor(.white)
                            
                            TextField("Enter area name", text: $areaName)
                                .textFieldStyle(CustomTextFieldStyle())
                        }
                        .padding(.horizontal, 20)
                        
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Description")
                                .font(.system(size: 18, weight: .semibold))
                                .foregroundColor(.white)
                            
                            TextField("Describe what this area does", text: $areaDescription, axis: .vertical)
                                .textFieldStyle(CustomTextFieldStyle())
                                .lineLimit(3...6)
                        }
                        .padding(.horizontal, 20)
                        
                        VStack(alignment: .leading, spacing: 16) {
                            HStack {
                                Text("Action Service")
                                    .font(.system(size: 18, weight: .semibold))
                                    .foregroundColor(.white)
                                
                                Spacer()
                                
                                if let service = selectedActionService {
                                    Button(action: {
                                        showingActionSelection = true
                                    }) {
                                        HStack(spacing: 8) {
                                            Image(systemName: service.icon)
                                                .foregroundColor(service.color)
                                            Text(service.name)
                                                .foregroundColor(.white)
                                            Image(systemName: "chevron.right")
                                                .foregroundColor(.gray)
                                                .font(.system(size: 12))
                                        }
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 8)
                                        .background(
                                            RoundedRectangle(cornerRadius: 8)
                                                .fill(Color.gray.opacity(0.2))
                                        )
                                    }
                                } else {
                                    Button(action: {
                                        showingActionSelection = true
                                    }) {
                                        HStack(spacing: 8) {
                                            Image(systemName: "plus")
                                            Text("Select Action Service")
                                        }
                                        .foregroundColor(AppColors.primaryBlue)
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 8)
                                        .background(
                                            RoundedRectangle(cornerRadius: 8)
                                                .stroke(AppColors.primaryBlue, lineWidth: 1)
                                        )
                                    }
                                }
                            }
                            
                            Text("The service that triggers the action")
                                .font(.system(size: 14))
                                .foregroundColor(.gray)
                        }
                        .padding(.horizontal, 20)
                        
                        Image(systemName: "arrow.down")
                            .font(.system(size: 24))
                            .foregroundColor(AppColors.primaryBlue)
                            .padding(.vertical, 10)
                        
                        VStack(alignment: .leading, spacing: 16) {
                            HStack {
                                Text("Reaction Service")
                                    .font(.system(size: 18, weight: .semibold))
                                    .foregroundColor(.white)
                                
                                Spacer()
                                
                                if let service = selectedReactionService {
                                    Button(action: {
                                        showingReactionSelection = true
                                    }) {
                                        HStack(spacing: 8) {
                                            Image(systemName: service.icon)
                                                .foregroundColor(service.color)
                                            Text(service.name)
                                                .foregroundColor(.white)
                                            Image(systemName: "chevron.right")
                                                .foregroundColor(.gray)
                                                .font(.system(size: 12))
                                        }
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 8)
                                        .background(
                                            RoundedRectangle(cornerRadius: 8)
                                                .fill(Color.gray.opacity(0.2))
                                        )
                                    }
                                } else {
                                    Button(action: {
                                        showingReactionSelection = true
                                    }) {
                                        HStack(spacing: 8) {
                                            Image(systemName: "plus")
                                            Text("Select Reaction Service")
                                        }
                                        .foregroundColor(AppColors.primaryBlue)
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 8)
                                        .background(
                                            RoundedRectangle(cornerRadius: 8)
                                                .stroke(AppColors.primaryBlue, lineWidth: 1)
                                        )
                                    }
                                }
                            }
                            
                            Text("The service that performs the reaction")
                                .font(.system(size: 14))
                                .foregroundColor(.gray)
                        }
                        .padding(.horizontal, 20)
                        
                        Button(action: createArea) {
                            HStack {
                                Image(systemName: "plus.circle.fill")
                                Text("Create AREA")
                            }
                            .font(.system(size: 18, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(
                                RoundedRectangle(cornerRadius: 12)
                                    .fill(canCreateArea ? AppGradients.button : LinearGradient(
                                        colors: [Color.gray.opacity(0.3)],
                                        startPoint: .leading,
                                        endPoint: .trailing
                                    ))
                            )
                        }
                        .disabled(!canCreateArea)
                        .padding(.horizontal, 20)
                        .padding(.bottom, 40)
                    }
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .navigationBarBackButtonHidden(true)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        presentationMode.wrappedValue.dismiss()
                    }
                    .foregroundColor(.white)
                }
            }
        }
        .sheet(isPresented: $showingActionSelection) {
            ServiceSelectionSheet(
                title: "Select Action Service",
                selectedService: $selectedActionService
            )
        }
        .sheet(isPresented: $showingReactionSelection) {
            ServiceSelectionSheet(
                title: "Select Reaction Service",
                selectedService: $selectedReactionService
            )
        }
        .alert("AREA Created!", isPresented: $showingSuccessAlert) {
            Button("OK") {
                presentationMode.wrappedValue.dismiss()
            }
        } message: {
            Text("Your AREA '\(areaName)' has been created successfully!")
        }
    }
    
    private var canCreateArea: Bool {
        !areaName.isEmpty && 
        !areaDescription.isEmpty && 
        selectedActionService != nil && 
        selectedReactionService != nil
    }
    
    private func createArea() {
        print("Creating AREA: \(areaName)")
        print("Description: \(areaDescription)")
        print("Action: \(selectedActionService?.name ?? "None")")
        print("Reaction: \(selectedReactionService?.name ?? "None")")
        
        showingSuccessAlert = true
    }
}

struct CustomTextFieldStyle: TextFieldStyle {
    func _body(configuration: TextField<Self._Label>) -> some View {
        configuration
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.gray.opacity(0.2))
            )
            .foregroundColor(.white)
    }
}

struct ServiceSelectionSheet: View {
    let title: String
    @Binding var selectedService: Service?
    @Environment(\.presentationMode) var presentationMode
    
    var body: some View {
        NavigationView {
            ZStack {
                Color.black
                    .ignoresSafeArea()
                
                ServiceSelectionView(
                    title: title,
                    selectedService: $selectedService,
                    onServiceSelected: { service in
                        selectedService = service
                        presentationMode.wrappedValue.dismiss()
                    }
                )
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        presentationMode.wrappedValue.dismiss()
                    }
                    .foregroundColor(.white)
                }
            }
        }
    }
}

#Preview {
    NewAreaView()
}
