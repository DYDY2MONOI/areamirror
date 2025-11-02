//
//  CreateView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct CreateView: View {
    @StateObject private var authService = AuthService.shared
    @StateObject private var catalogService = CatalogService.shared
    @StateObject private var areaService = AreaService.shared
    @State private var selectedAction: AboutAction?
    @State private var selectedReaction: AboutReaction?
    @State private var selectedActionService: AboutService?
    @State private var selectedReactionService: AboutService?
    @State private var areaName: String = ""
    @State private var areaDescription: String = ""
    @State private var showActionPicker = false
    @State private var showReactionPicker = false
    @State private var isCreating = false
    @State private var errorMessage: String?
    @State private var showSuccessAlert = false
    
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
                            if catalogService.isLoading {
                                ProgressView("Loading services...")
                                    .foregroundColor(.white)
                            } else if let error = catalogService.errorMessage {
                                Text(error)
                                    .foregroundColor(.red)
                            }
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
                                        Text(selectedAction?.name ?? "Select Action")
                                            .foregroundColor(selectedAction == nil ? .gray : .white)
                                        
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
                                        Text(selectedReaction?.name ?? "Select Reaction")
                                            .foregroundColor(selectedReaction == nil ? .gray : .white)
                                        
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
                                    .autocapitalization(.none)
                                    .disableAutocorrection(true)
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
                                    .autocapitalization(.none)
                                    .disableAutocorrection(true)
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
                            
                            if let errorMessage = errorMessage {
                                Text(errorMessage)
                                    .foregroundColor(.red)
                                    .font(.caption)
                                    .padding(.horizontal, 20)
                            }
                            
                            Button(action: {
                                createArea()
                            }) {
                                HStack {
                                    if isCreating {
                                        ProgressView()
                                            .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                    }
                                    Text("CREATE AREA")
                                        .font(.headline)
                                        .foregroundColor(.white)
                                }
                                .frame(maxWidth: .infinity)
                                .padding(.vertical, 16)
                                .background(buttonBackground)
                                .cornerRadius(12)
                            }
                            .disabled(!canCreate || isCreating)
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
            ActionPickerView(
                selectedAction: $selectedAction,
                selectedService: $selectedActionService
            )
            .environmentObject(catalogService)
        }
        .sheet(isPresented: $showReactionPicker) {
            ReactionPickerView(
                selectedReaction: $selectedReaction,
                selectedService: $selectedReactionService
            )
            .environmentObject(catalogService)
        }
        .alert("AREA Created!", isPresented: $showSuccessAlert) {
            Button("OK") {
                selectedAction = nil
                selectedReaction = nil
                selectedActionService = nil
                selectedReactionService = nil
                areaName = ""
                areaDescription = ""
            }
        } message: {
            Text("Your AREA '\(areaName)' has been created successfully!")
        }
        .task {
            if catalogService.services.isEmpty && !catalogService.isLoading {
                await catalogService.fetchCatalog()
            }
        }
    }
    
    private var canCreate: Bool {
        selectedAction != nil && selectedReaction != nil && !areaName.isEmpty
    }
    
    private var buttonBackground: some View {
        Group {
            if selectedAction == nil || selectedReaction == nil || areaName.isEmpty {
                LinearGradient(
                    gradient: Gradient(colors: [Color.gray, Color.gray]),
                    startPoint: .leading,
                    endPoint: .trailing
                )
            } else {
                AppGradients.button
            }
        }
    }
    
    private func createArea() {
        guard let actionService = selectedActionService?.name,
              let reactionService = selectedReactionService?.name,
              let triggerAction = selectedAction,
              selectedReaction != nil else {
            errorMessage = "Please select both an action and a reaction"
            return
        }

        isCreating = true
        errorMessage = nil

        Task {
            do {
                let resolvedTriggerType = AreaPayloadDefaults.triggerType(for: actionService)
                let resolvedActionType = AreaPayloadDefaults.actionType(for: reactionService)
                let triggerConfig = AreaPayloadDefaults.triggerConfig(for: actionService, actionName: triggerAction.name)
                let actionConfig = AreaPayloadDefaults.actionConfig(for: reactionService)

                let payload = AreaService.CreateOrUpdateAreaRequest(
                    name: areaName,
                    description: areaDescription.isEmpty ? nil : areaDescription,
                    triggerService: selectedActionService?.id ?? actionService,
                    triggerType: resolvedTriggerType,
                    actionService: selectedReactionService?.id ?? reactionService,
                    actionType: resolvedActionType,
                    triggerConfig: triggerConfig.isEmpty ? nil : triggerConfig,
                    actionConfig: actionConfig.isEmpty ? nil : actionConfig,
                    isActive: true
                )
                
                _ = try await areaService.createArea(payload: payload)
                print("✅ Area created successfully")
                
                DispatchQueue.main.async {
                    isCreating = false
                    showSuccessAlert = true
                }
            } catch {
                print("❌ Error creating area: \(error.localizedDescription)")
                DispatchQueue.main.async {
                    isCreating = false
                    errorMessage = error.localizedDescription
                }
            }
        }
    }
}

struct ActionPickerView: View {
    @EnvironmentObject var catalogService: CatalogService
    @Binding var selectedAction: AboutAction?
    @Binding var selectedService: AboutService?
    @Environment(\.dismiss) private var dismiss
    
    private var servicesWithActions: [AboutService] {
        catalogService.services.filter { !$0.actions.isEmpty }
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                if catalogService.isLoading {
                    VStack {
                        ProgressView("Loading actions...")
                            .foregroundColor(.white)
                    }
                } else if let error = catalogService.errorMessage {
                    VStack(spacing: 16) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.system(size: 50))
                            .foregroundColor(.red)
                        Text("Error loading actions")
                            .font(.headline)
                            .foregroundColor(.white)
                        Text(error)
                            .font(.caption)
                            .foregroundColor(.gray)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal)
                    }
                } else if servicesWithActions.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "list.bullet.rectangle")
                            .font(.system(size: 50))
                            .foregroundColor(.gray)
                        Text("No actions available")
                            .font(.headline)
                            .foregroundColor(.white)
                        Text("No services with actions were found.")
                            .font(.caption)
                            .foregroundColor(.gray)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal)
                    }
                } else {
                    List {
                        ForEach(servicesWithActions) { service in
                            Section(service.name) {
                                ForEach(service.actions) { action in
                                    Button(action: {
                                        selectedAction = action
                                        selectedService = service
                                        dismiss()
                                    }) {
                                        HStack {
                                            ServiceIconView(name: service.name)
                                            VStack(alignment: .leading, spacing: 4) {
                                                Text(action.name)
                                                    .foregroundColor(.white)
                                                Text(action.description)
                                                    .font(.caption)
                                                    .foregroundColor(.gray)
                                            }
                                            Spacer()
                                            if selectedAction?.id == action.id {
                                                Image(systemName: "checkmark")
                                                    .foregroundColor(AppColors.primaryBlue)
                                            }
                                        }
                                        .padding(.vertical, 8)
                                    }
                                    .listRowBackground(Color.clear)
                                }
                            }
                        }
                    }
                    .listStyle(PlainListStyle())
                }
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
        .task {
            if catalogService.services.isEmpty && !catalogService.isLoading {
                await catalogService.fetchCatalog()
            }
        }
    }
}

struct ReactionPickerView: View {
    @EnvironmentObject var catalogService: CatalogService
    @Binding var selectedReaction: AboutReaction?
    @Binding var selectedService: AboutService?
    @Environment(\.dismiss) private var dismiss
    
    private var servicesWithReactions: [AboutService] {
        catalogService.services.filter { !$0.reactions.isEmpty }
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                if catalogService.isLoading {
                    VStack {
                        ProgressView("Loading reactions...")
                            .foregroundColor(.white)
                    }
                } else if let error = catalogService.errorMessage {
                    VStack(spacing: 16) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.system(size: 50))
                            .foregroundColor(.red)
                        Text("Error loading reactions")
                            .font(.headline)
                            .foregroundColor(.white)
                        Text(error)
                            .font(.caption)
                            .foregroundColor(.gray)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal)
                    }
                } else if servicesWithReactions.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "list.bullet.rectangle")
                            .font(.system(size: 50))
                            .foregroundColor(.gray)
                        Text("No reactions available")
                            .font(.headline)
                            .foregroundColor(.white)
                        Text("No services with reactions were found.")
                            .font(.caption)
                            .foregroundColor(.gray)
                            .multilineTextAlignment(.center)
                            .padding(.horizontal)
                    }
                } else {
                    List {
                        ForEach(servicesWithReactions) { service in
                            Section(service.name) {
                                ForEach(service.reactions) { reaction in
                                    Button(action: {
                                        selectedReaction = reaction
                                        selectedService = service
                                        dismiss()
                                    }) {
                                        HStack {
                                            ServiceIconView(name: service.name)
                                            VStack(alignment: .leading, spacing: 4) {
                                                Text(reaction.name)
                                                    .foregroundColor(.white)
                                                Text(reaction.description)
                                                    .font(.caption)
                                                    .foregroundColor(.gray)
                                            }
                                            Spacer()
                                            if selectedReaction?.id == reaction.id {
                                                Image(systemName: "checkmark")
                                                    .foregroundColor(AppColors.primaryBlue)
                                            }
                                        }
                                        .padding(.vertical, 8)
                                    }
                                    .listRowBackground(Color.clear)
                                }
                            }
                        }
                    }
                    .listStyle(PlainListStyle())
                }
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
        .task {
            if catalogService.services.isEmpty && !catalogService.isLoading {
                await catalogService.fetchCatalog()
            }
        }
    }
}

#Preview {
    CreateView()
}

private struct ServiceIconView: View {
    let name: String

    var body: some View {
        ZStack {
            Circle()
                .fill(iconBackgroundColor(for: name))
                .frame(width: 32, height: 32)
            Image(systemName: iconSymbol(for: name))
                .foregroundColor(.white)
                .font(.system(size: 16, weight: .bold))
        }
    }

    private func iconSymbol(for service: String) -> String {
        switch service.lowercased() {
        case "gmail": return "envelope.fill"
        case "slack": return "bubble.left.and.bubble.right.fill"
        case "github": return "chevron.left.slash.chevron.right"
        case "weather": return "cloud.sun.fill"
        case "google calendar": return "calendar"
        case "discord": return "bubble.left.fill"
        case "onedrive": return "icloud.and.arrow.up"
        case "google sheets": return "tablecells"
        case "spotify": return "music.note"
        default: return "bolt.fill"
        }
    }

    private func iconBackgroundColor(for service: String) -> Color {
        switch service.lowercased() {
        case "gmail": return Color(red: 0.92, green: 0.26, blue: 0.21)
        case "slack": return Color(red: 0.2, green: 0.4, blue: 0.8)
        case "github": return Color(red: 0.16, green: 0.16, blue: 0.16)
        case "weather": return Color(red: 0.0, green: 0.5, blue: 1.0)
        case "google calendar": return Color(red: 0.16, green: 0.47, blue: 0.95)
        case "discord": return Color(red: 0.35, green: 0.4, blue: 0.95)
        case "onedrive": return Color(red: 0.0, green: 0.5, blue: 0.9)
        case "google sheets": return Color(red: 0.0, green: 0.6, blue: 0.2)
        case "spotify": return Color(red: 0.2, green: 0.8, blue: 0.2)
        default: return AppColors.primaryBlue
        }
    }
}
