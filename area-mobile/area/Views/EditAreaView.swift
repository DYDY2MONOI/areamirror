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
                AppGradients.background.ignoresSafeArea()
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        headerCard
                        
                        formCard
                        
                        if template.triggerService == "Google Calendar" {
                            calendarTriggerCard
                        }
                        
                        if template.actionService == "Gmail" {
                            gmailActionCard
                        }
                        
                        if let errorMessage = errorMessage {
                            Text(errorMessage)
                                .foregroundColor(.red)
                                .padding(.horizontal, 20)
                        }
                        
                        Button(action: saveArea) {
                            HStack(spacing: 10) {
                                if isSaving { ProgressView().tint(.white) }
                                Image(systemName: existingArea == nil ? "tray.and.arrow.down.fill" : "arrow.triangle.2.circlepath.circle.fill")
                                    .foregroundColor(.white)
                                Text(existingArea == nil ? "Save AREA" : "Update AREA")
                            }
                            .font(.system(size: 18, weight: .semibold))
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(
                                RoundedRectangle(cornerRadius: 14)
                                    .fill(canSave ? AppGradients.button : LinearGradient(colors: [Color.gray.opacity(0.3)], startPoint: .leading, endPoint: .trailing))
                            )
                        }
                        .disabled(!canSave || isSaving)
                        .padding(.horizontal, 20)
                        .padding(.bottom, 40)
                    }
                    .padding(.top, 16)
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: { presentationMode.wrappedValue.dismiss() }) {
                        Image(systemName: "xmark")
                            .font(.system(size: 16, weight: .bold))
                            .foregroundColor(.white)
                            .padding(8)
                            .background(Color.white.opacity(0.08))
                            .clipShape(Circle())
                    }
                    .buttonStyle(PlainButtonStyle())
                }
            }
            .toolbarBackground(.hidden, for: .navigationBar)
            .toolbarColorScheme(.dark, for: .navigationBar)
            .navigationBarBackButtonHidden(true)
        }
        .onAppear(perform: initializeForm)
    }
    
    private var calendarTriggerCard: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: "calendar")
                    .foregroundColor(.white)
                Text("Calendar Trigger")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(.white)
            }
            
            labeledField(label: "Event date", placeholder: "YYYY-MM-DD", text: $eventDate)
            labeledField(label: "Event time", placeholder: "HH:MM", text: $eventTime)
            labeledField(label: "Event title", placeholder: "Event title", text: $eventTitle)
            labeledField(label: "Calendar ID", placeholder: "primary", text: $calendarId)
        }
        .cardStyle()
        .padding(.horizontal, 20)
    }
    
    private var gmailActionCard: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: "envelope.fill")
                    .foregroundColor(.white)
                Text("Gmail Action")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(.white)
            }
            
            labeledField(label: "To email", placeholder: "your-email@gmail.com", text: $toEmail)
            labeledField(label: "Subject", placeholder: "Reminder: {{eventTitle}}", text: $subject)
            labeledField(label: "Body", placeholder: "Body", text: $emailBody, axis: .vertical)
        }
        .cardStyle()
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

// MARK: - Styled Sections
extension EditAreaView {
    private var headerCard: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text(existingArea == nil ? "Create AREA" : "Edit AREA")
                .font(.system(size: 28, weight: .bold))
                .foregroundColor(.white)
            
            ZStack {
                RoundedRectangle(cornerRadius: 16)
                    .fill(LinearGradient(
                        gradient: Gradient(colors: [serviceColor(template.triggerService), serviceColor(template.actionService)]),
                        startPoint: .topLeading,
                        endPoint: .bottomTrailing
                    ))
                    .frame(maxWidth: .infinity)
                    .frame(height: 140)
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.white.opacity(0.1), lineWidth: 1)
                    )
                
                VStack {
                    HStack(alignment: .center) {
                        VStack(alignment: .leading, spacing: 6) {
                            Text(name.isEmpty ? template.title : name)
                                .font(.system(size: 20, weight: .semibold))
                                .foregroundColor(.white)
                                .shadow(color: .black.opacity(0.6), radius: 6, x: 0, y: 2)
                            Text(template.subtitle.isEmpty ? "Calendar automation" : template.subtitle)
                                .font(.system(size: 14))
                                .foregroundColor(.white.opacity(0.9))
                                .shadow(color: .black.opacity(0.5), radius: 6, x: 0, y: 2)
                        }
                        Spacer()
                        ZStack {
                            Circle()
                                .fill(Color.white.opacity(0.12))
                                .frame(width: 48, height: 48)
                            Image(systemName: serviceIcon(template.actionService))
                                .font(.system(size: 22, weight: .bold))
                                .foregroundColor(.white)
                                .shadow(color: .black.opacity(0.6), radius: 6, x: 0, y: 2)
                        }
                    }
                    .padding(16)
                    Spacer()
                }
            }
        }
        .padding(.horizontal, 20)
    }
    
    private var formCard: some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: "slider.horizontal.3")
                    .foregroundColor(.white)
                Text("Details")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(.white)
            }
            labeledField(label: "Name", placeholder: "Area name", text: $name)
            labeledField(label: "Description", placeholder: "Description", text: $description)
        }
        .cardStyle()
        .padding(.horizontal, 20)
    }
    
    @ViewBuilder
    private func labeledField(label: String, placeholder: String, text: Binding<String>, axis: Axis.Set = []) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(label)
                .foregroundColor(.white)
                .font(.system(size: 14, weight: .medium))
            
            if axis == .vertical {
                TextField(placeholder, text: text, axis: .vertical)
                    .textFieldStyle(CustomTextFieldStyle())
                    .lineLimit(3...6)
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            } else {
                TextField(placeholder, text: text)
                    .textFieldStyle(CustomTextFieldStyle())
                    .autocapitalization(.none)
                    .disableAutocorrection(true)
            }
        }
    }
}

// MARK: - Modifiers & Helpers
fileprivate struct CardModifier: ViewModifier {
    func body(content: Content) -> some View {
        content
            .padding(16)
            .background(
                RoundedRectangle(cornerRadius: 16)
                    .fill(AppColors.darkBackground)
                    .overlay(
                        RoundedRectangle(cornerRadius: 16)
                            .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                    )
            )
    }
}

fileprivate extension View {
    func cardStyle() -> some View { self.modifier(CardModifier()) }
}

fileprivate func serviceIcon(_ service: String) -> String {
    switch service.lowercased() {
    case "github": return "hammer.fill"
    case "gmail": return "envelope.fill"
    case "discord": return "message.fill"
    case "slack": return "message.circle.fill"
    case "weather": return "cloud.sun.fill"
    case "instagram": return "camera.fill"
    case "twitter": return "bird.fill"
    case "youtube": return "play.rectangle.fill"
    case "spotify": return "music.note"
    case "telegram": return "paperplane.fill"
    case "dropbox": return "folder.fill"
    case "notion": return "doc.text.fill"
    default: return "gearshape.fill"
    }
}

fileprivate func serviceColor(_ service: String) -> Color {
    switch service.lowercased() {
    case "github": return Color(red: 0.2, green: 0.2, blue: 0.2)
    case "gmail": return Color(red: 0.92, green: 0.26, blue: 0.21)
    case "discord": return Color(red: 0.35, green: 0.4, blue: 0.95)
    case "slack": return Color(red: 0.36, green: 0.8, blue: 0.36)
    case "weather": return Color(red: 0.0, green: 0.7, blue: 1.0)
    case "instagram": return Color(red: 0.8, green: 0.2, blue: 0.6)
    case "twitter": return Color(red: 0.1, green: 0.6, blue: 0.9)
    case "youtube": return Color(red: 1.0, green: 0.0, blue: 0.0)
    case "spotify": return Color(red: 0.2, green: 0.8, blue: 0.2)
    case "telegram": return Color(red: 0.0, green: 0.7, blue: 0.9)
    case "dropbox": return Color(red: 0.0, green: 0.5, blue: 0.8)
    case "notion": return Color(red: 0.2, green: 0.2, blue: 0.2)
    default: return AppColors.primaryBlue
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
