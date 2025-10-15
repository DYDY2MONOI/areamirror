//
//  EditAreaView.swift
//  area
//
//  Created by Assistant on 15/10/2025.
//

import SwiftUI

struct EditAreaView: View {
    @Environment(\.presentationMode) var presentationMode
    @StateObject private var areaService = AreaService.shared
    
    let template: AreaTemplate
    let existingArea: Area?
    
    @State private var name: String = ""
    @State private var description: String = ""
    
    // Trigger (calendar/github/weather minimal fields supported now)
    @State private var eventDate: String = ""
    @State private var eventTime: String = ""
    @State private var eventTitle: String = ""
    @State private var calendarId: String = "primary"
    
    // Action (gmail minimal fields supported now)
    @State private var toEmail: String = ""
    @State private var subject: String = ""
    @State private var emailBody: String = ""
    
    @State private var isSaving = false
    @State private var errorMessage: String?
    
    init(template: AreaTemplate, existingArea: Area? = nil) {
        self.template = template
        self.existingArea = existingArea
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                Color.black.ignoresSafeArea()
                ScrollView {
                    VStack(alignment: .leading, spacing: 20) {
                        Group {
                            Text(existingArea == nil ? "Create AREA" : "Edit AREA")
                                .font(.system(size: 28, weight: .bold))
                                .foregroundColor(.white)
                            Text("\(template.triggerService) → \(template.actionService)")
                                .foregroundColor(.gray)
                        }
                        .padding(.horizontal, 20)
                        
                        VStack(alignment: .leading, spacing: 12) {
                            Text("Name").foregroundColor(.white)
                            TextField("Area name", text: $name)
                                .textFieldStyle(CustomTextFieldStyle())
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                            Text("Description").foregroundColor(.white)
                            TextField("Description", text: $description)
                                .textFieldStyle(CustomTextFieldStyle())
                                .autocapitalization(.none)
                                .disableAutocorrection(true)
                        }
                        .padding(.horizontal, 20)
                        
                        if template.triggerService == "Google Calendar" {
                            calendarTriggerConfig
                        }
                        
                        if template.actionService == "Gmail" {
                            gmailActionConfig
                        }
                        
                        if let errorMessage = errorMessage {
                            Text(errorMessage)
                                .foregroundColor(.red)
                                .padding(.horizontal, 20)
                        }
                        
                        Button(action: saveArea) {
                            HStack {
                                if isSaving { ProgressView().tint(.white) }
                                Text(existingArea == nil ? "Save AREA" : "Update AREA")
                            }
                            .font(.system(size: 18, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(
                                RoundedRectangle(cornerRadius: 12)
                                    .fill(canSave ? AppGradients.button : LinearGradient(colors: [Color.gray.opacity(0.3)], startPoint: .leading, endPoint: .trailing))
                            )
                        }
                        .disabled(!canSave || isSaving)
                        .padding(.horizontal, 20)
                        .padding(.bottom, 40)
                    }
                    .padding(.top, 20)
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("Close") { presentationMode.wrappedValue.dismiss() }
                        .foregroundColor(.white)
                }
            }
        }
        .onAppear(perform: initializeForm)
    }
    
    private var calendarTriggerConfig: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Calendar Trigger").font(.system(size: 18, weight: .semibold)).foregroundColor(.white)
            Text("Event date").foregroundColor(.white)
            TextField("YYYY-MM-DD", text: $eventDate).textFieldStyle(CustomTextFieldStyle())
            Text("Event time").foregroundColor(.white)
            TextField("HH:MM", text: $eventTime).textFieldStyle(CustomTextFieldStyle())
            Text("Event title").foregroundColor(.white)
            TextField("Event title", text: $eventTitle).textFieldStyle(CustomTextFieldStyle())
            Text("Calendar ID").foregroundColor(.white)
            TextField("primary", text: $calendarId).textFieldStyle(CustomTextFieldStyle())
        }
        .padding(.horizontal, 20)
    }
    
    private var gmailActionConfig: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Gmail Action").font(.system(size: 18, weight: .semibold)).foregroundColor(.white)
            Text("To email").foregroundColor(.white)
            TextField("your-email@gmail.com", text: $toEmail).textFieldStyle(CustomTextFieldStyle())
            Text("Subject").foregroundColor(.white)
            TextField("Reminder: {{eventTitle}}", text: $subject).textFieldStyle(CustomTextFieldStyle())
            Text("Body").foregroundColor(.white)
            TextField("Body", text: $emailBody, axis: .vertical).textFieldStyle(CustomTextFieldStyle()).lineLimit(3...6)
        }
        .padding(.horizontal, 20)
    }
    
    private var canSave: Bool {
        !name.isEmpty && (template.actionService != "Gmail" || (!toEmail.isEmpty && !subject.isEmpty))
    }
    
    private func initializeForm() {
        if let existing = existingArea {
            name = existing.name
            description = existing.description
        } else {
            name = template.title
            description = template.description
        }
        
        if template.actionService == "Gmail" {
            if subject.isEmpty { subject = "Reminder: {{eventTitle}}" }
            if emailBody.isEmpty { emailBody = "Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}" }
        }
    }
    
    private func saveArea() {
        isSaving = true
        errorMessage = nil
        
        Task {
            do {
                var triggerConfig: [String: AnyCodable] = [:]
                if template.triggerService == "Google Calendar" {
                    let combinedTime = (!eventDate.isEmpty && !eventTime.isEmpty) ? "\(eventDate)T\(eventTime):00Z" : eventTime
                    triggerConfig = [
                        "eventDate": AnyCodable(eventDate),
                        "eventTime": AnyCodable(combinedTime),
                        "eventTitle": AnyCodable(eventTitle),
                        "calendarId": AnyCodable(calendarId)
                    ]
                }
                
                var actionConfig: [String: AnyCodable] = [:]
                if template.actionService == "Gmail" {
                    actionConfig = [
                        "toEmail": AnyCodable(toEmail),
                        "subject": AnyCodable(subject),
                        "body": AnyCodable(emailBody)
                    ]
                }
                
                let payload = AreaService.CreateOrUpdateAreaRequest(
                    name: name,
                    description: description,
                    triggerService: template.triggerService,
                    triggerType: template.triggerService == "Google Calendar" ? "Event" : "Webhook",
                    actionService: template.actionService,
                    actionType: resolveActionType(template.actionService),
                    triggerConfig: triggerConfig,
                    actionConfig: actionConfig
                )
                
                if let area = existingArea {
                    _ = try await areaService.updateArea(areaId: area.id, payload: payload)
                } else {
                    _ = try await areaService.createArea(payload: payload)
                }
                
                DispatchQueue.main.async {
                    isSaving = false
                    presentationMode.wrappedValue.dismiss()
                }
            } catch {
                DispatchQueue.main.async {
                    isSaving = false
                    errorMessage = error.localizedDescription
                }
            }
        }
    }
    
    private func resolveActionType(_ service: String) -> String {
        switch service {
        case "Gmail": return "SendEmail"
        case "Discord": return "SendDiscordMessage"
        default: return "Action"
        }
    }
}

// MARK: - Preview
#Preview {
    EditAreaView(
        template: AreaTemplate(
            id: "1",
            title: "Google Calendar → Gmail",
            subtitle: "Send yourself an email when a calendar event starts",
            description: "This AREA sends you an email via Gmail when a new event from your Google Calendar starts.",
            icon: "calendar",
            gradientClass: "blue",
            triggerService: "Google Calendar",
            actionService: "Gmail",
            triggerIconURL: nil,
            actionIconURL: nil,
            isActive: true
        )
    )
}
