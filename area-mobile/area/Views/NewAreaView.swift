//
//  NewAreaView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct NewAreaView: View {
    @Environment(\.presentationMode) var presentationMode
    @StateObject private var areaService = AreaService.shared
    @State private var areaName = ""
    @State private var areaDescription = ""
    @State private var selectedActions: [Service] = []
    @State private var selectedReactions: [Service] = []
    @State private var showingActionSelection = false
    @State private var showingReactionSelection = false
    @State private var showingSuccessAlert = false
    @State private var isCreating = false
    @State private var errorMessage: String?
    
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
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                        }
                        .padding(.horizontal, 20)
                        
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Description")
                                .font(.system(size: 18, weight: .semibold))
                                .foregroundColor(.white)
                            
                            TextField("Describe what this area does", text: $areaDescription, axis: .vertical)
                                .textFieldStyle(CustomTextFieldStyle())
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                                .lineLimit(3...6)
                        }
                        .padding(.horizontal, 20)
                        
                        VStack(alignment: .leading, spacing: 16) {
                            HStack {
                                Text("Action Services")
                                    .font(.system(size: 18, weight: .semibold))
                                    .foregroundColor(.white)
                                
                                Spacer()
                                
                                Button(action: {
                                    showingActionSelection = true
                                }) {
                                    HStack(spacing: 8) {
                                        Image(systemName: "plus")
                                        Text("Add Action")
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
                            
                            Text("Services that trigger actions")
                                .font(.system(size: 14))
                                .foregroundColor(.gray)
                            
                            if selectedActions.isEmpty {
                                Text("No actions selected")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                                    .italic()
                                    .padding(.vertical, 8)
                            } else {
                                ForEach(Array(selectedActions.enumerated()), id: \.offset) { index, service in
                                    HStack {
                                        ZStack {
                                            RoundedRectangle(cornerRadius: 3)
                                                .fill(service.color == .white ? Color.gray : service.color)
                                                .frame(width: 16, height: 16)
                                            
                                            Image(service.icon)
                                                .resizable()
                                                .aspectRatio(contentMode: .fit)
                                                .frame(width: 14, height: 14)
                                                .colorMultiply(service.color == .white ? .white : service.color)
                                        }
                                        
                                        Text(service.name)
                                            .foregroundColor(.white)
                                        
                                        Spacer()
                                        
                                        Button(action: {
                                            selectedActions.remove(at: index)
                                        }) {
                                            Image(systemName: "xmark.circle.fill")
                                                .foregroundColor(.red)
                                                .font(.system(size: 16))
                                        }
                                    }
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 8)
                                    .background(
                                        RoundedRectangle(cornerRadius: 8)
                                            .fill(Color.gray.opacity(0.2))
                                    )
                                }
                            }
                        }
                        .padding(.horizontal, 20)
                        
                        Image(systemName: "arrow.down")
                            .font(.system(size: 24))
                            .foregroundColor(AppColors.primaryBlue)
                            .padding(.vertical, 10)
                        
                        VStack(alignment: .leading, spacing: 16) {
                            HStack {
                                Text("Reaction Services")
                                    .font(.system(size: 18, weight: .semibold))
                                    .foregroundColor(.white)
                                
                                Spacer()
                                
                                Button(action: {
                                    showingReactionSelection = true
                                }) {
                                    HStack(spacing: 8) {
                                        Image(systemName: "plus")
                                        Text("Add Reaction")
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
                            
                            Text("Services that perform reactions")
                                .font(.system(size: 14))
                                .foregroundColor(.gray)
                            
                            if selectedReactions.isEmpty {
                                Text("No reactions selected")
                                    .font(.system(size: 14))
                                    .foregroundColor(.gray)
                                    .italic()
                                    .padding(.vertical, 8)
                            } else {
                                ForEach(Array(selectedReactions.enumerated()), id: \.offset) { index, service in
                                    HStack {
                                        ZStack {
                                            RoundedRectangle(cornerRadius: 3)
                                                .fill(service.color == .white ? Color.gray : service.color)
                                                .frame(width: 16, height: 16)
                                            
                                            Image(service.icon)
                                                .resizable()
                                                .aspectRatio(contentMode: .fit)
                                                .frame(width: 14, height: 14)
                                                .colorMultiply(service.color == .white ? .white : service.color)
                                        }
                                        
                                        Text(service.name)
                                            .foregroundColor(.white)
                                        
                                        Spacer()
                                        
                                        Button(action: {
                                            selectedReactions.remove(at: index)
                                        }) {
                                            Image(systemName: "xmark.circle.fill")
                                                .foregroundColor(.red)
                                                .font(.system(size: 16))
                                        }
                                    }
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 8)
                                    .background(
                                        RoundedRectangle(cornerRadius: 8)
                                            .fill(Color.gray.opacity(0.2))
                                    )
                                }
                            }
                        }
                        .padding(.horizontal, 20)
                        
                        Button(action: createArea) {
                            HStack(spacing: 8) {
                                if isCreating {
                                    ProgressView()
                                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                } else {
                                    Image(systemName: "plus.circle.fill")
                                }
                                Text(isCreating ? "Creating..." : "Create AREA")
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
                        .disabled(!canCreateArea || isCreating)
                        .padding(.horizontal, 20)
                        .padding(.bottom, 16)
                        
                        if let errorMessage = errorMessage {
                            Text(errorMessage)
                                .font(.system(size: 14))
                                .foregroundColor(.red)
                                .multilineTextAlignment(.center)
                                .padding(.horizontal, 24)
                        }
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
            MultiServiceSelectionSheet(
                title: "Select Action Services",
                mode: .action,
                selectedServices: $selectedActions
            )
            .environmentObject(CatalogService.shared)
        }
        .sheet(isPresented: $showingReactionSelection) {
            MultiServiceSelectionSheet(
                title: "Select Reaction Services",
                mode: .reaction,
                selectedServices: $selectedReactions
            )
            .environmentObject(CatalogService.shared)
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
        !selectedActions.isEmpty && 
        !selectedReactions.isEmpty
    }
    
    private func createArea() {
        guard let triggerService = selectedActions.first?.name,
              let actionService = selectedReactions.first?.name else {
            errorMessage = "Please select both an action and a reaction service"
            return
        }

        isCreating = true
        errorMessage = nil

        print("NewAreaView.createArea attempting with trigger=\(triggerService) -> action=\(actionService)")

        Task {
            do {
                let triggerType = AreaPayloadDefaults.triggerType(for: triggerService)
                let actionType = AreaPayloadDefaults.actionType(for: actionService)
                let triggerConfig = AreaPayloadDefaults.triggerConfig(for: triggerService, actionName: areaName.isEmpty ? triggerService : areaName)
                let actionConfig = AreaPayloadDefaults.actionConfig(for: actionService)

                let payload = AreaService.CreateOrUpdateAreaRequest(
                    name: areaName,
                    description: areaDescription.isEmpty ? nil : areaDescription,
                    triggerService: triggerService,
                    triggerType: triggerType,
                    actionService: actionService,
                    actionType: actionType,
                    triggerConfig: triggerConfig.isEmpty ? nil : triggerConfig,
                    actionConfig: actionConfig.isEmpty ? nil : actionConfig,
                    isActive: true
                )

                _ = try await areaService.createArea(payload: payload)
                print("Area created successfully from NewAreaView")

                await MainActor.run {
                    isCreating = false
                    showingSuccessAlert = true
                }
            } catch {
                print("Failed to create area from NewAreaView: \(error.localizedDescription)")
                await MainActor.run {
                    isCreating = false
                    errorMessage = error.localizedDescription
                }
            }
        }
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

struct MultiServiceSelectionSheet: View {
    enum Mode { case action, reaction }
    let title: String
    let mode: Mode
    @Binding var selectedServices: [Service]
    @Environment(\.presentationMode) var presentationMode
    @EnvironmentObject var catalogService: CatalogService
    
    private var availableServices: [Service] {
        let source = catalogService.services
        let filtered: [AboutService]
        switch mode {
        case .action:
            filtered = source.filter { !$0.actions.isEmpty }
        case .reaction:
            filtered = source.filter { !$0.reactions.isEmpty }
        }
        return filtered.map { svc in
            let iconName = iconForService(svc.name)
            return Service(name: svc.name, icon: iconName, color: Color.gray, category: .other)
        }
    }
    
    private func iconForService(_ name: String) -> String {
        switch name.lowercased() {
        case "gmail": return "gmail"
        case "google calendar": return "google-calendar"
        case "google sheets": return "google-sheets"
        case "github": return "github"
        case "discord": return "discord"
        case "slack": return "slack"
        case "weather": return "weather"
        case "onedrive": return "onedrive"
        case "spotify": return "spotify"
        default: return "generic-service"
        }
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                Color.black.ignoresSafeArea()
                
                ScrollView {
                    LazyVStack(spacing: 12) {
                        ForEach(availableServices, id: \.id) { service in
                            ServiceRow(
                                service: service,
                                isSelected: selectedServices.contains { $0.name == service.name },
                                onToggle: {
                                    if selectedServices.contains(where: { $0.name == service.name }) {
                                        selectedServices.removeAll { $0.name == service.name }
                                    } else {
                                        selectedServices = [service]
                                    }
                                }
                            )
                        }
                    }
                    .padding(.horizontal, 20)
                    .padding(.top, 20)
                }
            }
            .navigationTitle(title)
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Cancel") {
                        presentationMode.wrappedValue.dismiss()
                    }
                    .foregroundColor(.white)
                }
                
                ToolbarItem(placement: .navigationBarTrailing) {
                    Button("Done") {
                        presentationMode.wrappedValue.dismiss()
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

struct ServiceRow: View {
    let service: Service
    let isSelected: Bool
    let onToggle: () -> Void
    
    var body: some View {
        Button(action: onToggle) {
            HStack(spacing: 16) {
                ZStack {
                    RoundedRectangle(cornerRadius: 8)
                        .fill(service.color == .white ? Color.gray : service.color)
                        .frame(width: 40, height: 40)
                    
                    Image(service.icon)
                        .resizable()
                        .aspectRatio(contentMode: .fit)
                        .frame(width: 20, height: 20)
                        .colorMultiply(.white)
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text(service.name)
                        .font(.system(size: 16, weight: .medium))
                        .foregroundColor(.white)
                    
                    Text("Service description")
                        .font(.system(size: 14))
                        .foregroundColor(.gray)
                }
                
                Spacer()
                
                Image(systemName: isSelected ? "checkmark.circle.fill" : "circle")
                    .foregroundColor(isSelected ? AppColors.primaryBlue : .gray)
                    .font(.system(size: 24))
            }
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color.gray.opacity(0.1))
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(isSelected ? AppColors.primaryBlue : Color.clear, lineWidth: 2)
                    )
            )
        }
        .buttonStyle(PlainButtonStyle())
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
