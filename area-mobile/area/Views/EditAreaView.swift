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
    
    let area: Area
    let isEditing: Bool
    
    @State private var name: String = ""
    @State private var description: String = ""
    
    @State private var eventDate: String = ""
    @State private var eventTime: String = ""
    @State private var eventTitle: String = ""
    @State private var calendarId: String = "primary"
    @State private var eventDateTime: Date = Date()
    
    @State private var toEmail: String = ""
    @State private var subject: String = ""
    @State private var emailBody: String = ""
    
    @State private var isSaving = false
    @State private var errorMessage: String?
    
    init(area: Area, isEditing: Bool = true) {
        print("📱 EditAreaView init with area: \(area.name) (ID: \(area.id)), isEditing: \(isEditing)")
        self.area = area
        self.isEditing = isEditing
        
        _name = State(initialValue: area.name)
        _description = State(initialValue: area.description)
        
        if let triggerConfig = area.triggerConfig {
            print("📱 Loading trigger config for existing area: \(area.name)")
            if area.triggerService == "Google Calendar" {
                _eventDate = State(initialValue: triggerConfig["eventDate"]?.value as? String ?? "")
                _eventTime = State(initialValue: triggerConfig["eventTime"]?.value as? String ?? "")
                _eventTitle = State(initialValue: triggerConfig["eventTitle"]?.value as? String ?? "")
                _calendarId = State(initialValue: triggerConfig["calendarId"]?.value as? String ?? "primary")
                if let parsed = Self.parseCalendarDateTime(dateString: _eventDate.wrappedValue, timeString: _eventTime.wrappedValue) {
                    _eventDateTime = State(initialValue: parsed)
                }
                print("📱 Loaded calendar config: date=\(triggerConfig["eventDate"]?.value as? String ?? "nil"), time=\(triggerConfig["eventTime"]?.value as? String ?? "nil")")
            }
        } else if area.triggerService == "Google Calendar" {
            print("📱 No trigger config, using default calendar config")
            _eventDate = State(initialValue: "")
            _eventTime = State(initialValue: "")
            _eventDateTime = State(initialValue: Date())
            _eventTitle = State(initialValue: "")
            _calendarId = State(initialValue: "primary")
        }
        
        if let actionConfig = area.actionConfig {
            print("📱 Loading action config for existing area: \(area.name)")
            if area.actionService == "Gmail" {
                _toEmail = State(initialValue: actionConfig["toEmail"]?.value as? String ?? "")
                _subject = State(initialValue: actionConfig["subject"]?.value as? String ?? "Reminder: {{eventTitle}}")
                _emailBody = State(initialValue: actionConfig["body"]?.value as? String ?? "Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}")
                print("📱 Loaded Gmail config: to=\(actionConfig["toEmail"]?.value as? String ?? "nil"), subject=\(actionConfig["subject"]?.value as? String ?? "nil")")
            }
        } else if area.actionService == "Gmail" {
            print("📱 No action config, using default Gmail config")
            _toEmail = State(initialValue: "")
            _subject = State(initialValue: "Reminder: {{eventTitle}}")
            _emailBody = State(initialValue: "Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}")
        }
    }
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background.ignoresSafeArea()
                ScrollView(showsIndicators: false) {
                    VStack(spacing: 20) {
                        headerCard
                        
                        formCard
                        
                        if area.triggerService == "Google Calendar" {
                            calendarTriggerCard
                        }
                        
                        if area.actionService == "Gmail" {
                            gmailActionCard
                        }
                        
                        if let errorMessage = errorMessage {
                            Text(errorMessage)
                                .foregroundColor(.red)
                                .padding(.horizontal, 20)
                        }
                        
                        #if DEBUG
                        if isEditing {
                            VStack(spacing: 8) {
                                Text("DEBUG: Edit Mode")
                                    .foregroundColor(.green)
                                    .font(.caption)
                                Text("Area ID: \(area.id)")
                                    .foregroundColor(.green)
                                    .font(.caption)
                                Text("Will call: PUT /areas/\(area.id)")
                                    .foregroundColor(.green)
                                    .font(.caption)
                            }
                            .padding(.horizontal, 20)
                        }
                        #endif
                        
                        Button(action: saveArea) {
                            HStack(spacing: 10) {
                                if isSaving { ProgressView().tint(.white) }
                                Image(systemName: isEditing ? "arrow.triangle.2.circlepath.circle.fill" : "tray.and.arrow.down.fill")
                                    .foregroundColor(.white)
                                Text(isEditing ? "Update AREA" : "Save AREA")
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
            
            VStack(alignment: .leading, spacing: 8) {
                Text("Date & time")
                    .foregroundColor(.white)
                    .font(.system(size: 14, weight: .medium))
                DatePicker("", selection: $eventDateTime, displayedComponents: [.date, .hourAndMinute])
                    .labelsHidden()
                    .datePickerStyle(.compact)
                    .environment(\.timeZone, TimeZone(secondsFromGMT: 0)!)
                    .onChange(of: eventDateTime) { _ in
                        let (dateStr, timeStr) = Self.formatCalendarStrings(from: eventDateTime)
                        eventDate = dateStr
                        eventTime = timeStr
                    }
            }
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
        !name.isEmpty && (area.actionService != "Gmail" || (!toEmail.isEmpty && !subject.isEmpty))
    }
    
    private func saveArea() {
        isSaving = true
        errorMessage = nil
        
        Task {
            do {
                var triggerConfig: [String: AnyCodable] = [:]
                if area.triggerService == "Google Calendar" {
                    // eventTime already formatted as RFC3339 from DatePicker helper
                    triggerConfig = [
                        "eventDate": AnyCodable(eventDate),
                        "eventTime": AnyCodable(eventTime),
                        "eventTitle": AnyCodable(eventTitle),
                        "calendarId": AnyCodable(calendarId)
                    ]
                }
                
                var actionConfig: [String: AnyCodable] = [:]
                if area.actionService == "Gmail" {
                    actionConfig = [
                        "toEmail": AnyCodable(toEmail),
                        "subject": AnyCodable(subject),
                        "body": AnyCodable(emailBody)
                    ]
                }
                
                let fullPayload = AreaService.CreateOrUpdateAreaRequest(
                    name: name,
                    description: description,
                    triggerService: area.triggerService,
                    triggerType: area.triggerService == "Google Calendar" ? "Event" : "Webhook",
                    actionService: area.actionService,
                    actionType: resolveActionType(area.actionService),
                    triggerConfig: triggerConfig,
                    actionConfig: actionConfig,
                    isActive: true
                )
                
                if isEditing {
                    print("🔄 Updating existing area with ID: \(area.id)")
                    _ = try await areaService.updateArea(areaId: area.id, payload: fullPayload)
                    print("✅ Area updated successfully")
                } else {
                    print("➕ Creating new area")
                    _ = try await areaService.createArea(payload: fullPayload)
                    print("✅ Area created successfully")
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

extension EditAreaView {
    private static func parseCalendarDateTime(dateString: String, timeString: String) -> Date? {
        if !timeString.isEmpty {
            let iso1 = ISO8601DateFormatter()
            if let d = iso1.date(from: timeString) { return d }

            if !dateString.isEmpty {
                let cleanedTime = timeString.count == 5 ? "\(timeString):00" : timeString
                let combined = "\(dateString)T\(cleanedTime)Z"
                if let d2 = iso1.date(from: combined) { return d2 }
            }
        }
        if !dateString.isEmpty {
            let df = DateFormatter()
            df.locale = Locale(identifier: "en_US_POSIX")
            df.timeZone = TimeZone(secondsFromGMT: 0)
            df.dateFormat = "yyyy-MM-dd"
            return df.date(from: dateString)
        }
        return nil
    }

    private static func formatCalendarStrings(from date: Date) -> (String, String) {
        let df = DateFormatter()
        df.locale = Locale(identifier: "en_US_POSIX")
        df.timeZone = TimeZone(secondsFromGMT: 0)
        df.dateFormat = "yyyy-MM-dd"
        let dateStr = df.string(from: date)

        let iso = ISO8601DateFormatter()
        iso.timeZone = TimeZone(secondsFromGMT: 0)
        iso.formatOptions = [.withInternetDateTime, .withColonSeparatorInTime]
        let isoStr = iso.string(from: date)
        return (dateStr, isoStr)
    }
    private var headerCard: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack {
                Text(isEditing ? "Edit AREA" : "Create AREA")
                    .font(.system(size: 28, weight: .bold))
                    .foregroundColor(.white)
                
                if isEditing {
                    Text("(ID: \(area.id))")
                        .font(.system(size: 14, weight: .medium))
                        .foregroundColor(.white.opacity(0.7))
                        .padding(.leading, 8)
                }
            }
            
            ZStack {
                RoundedRectangle(cornerRadius: 16)
                    .fill(LinearGradient(
                        gradient: Gradient(colors: [serviceColor(area.triggerService), serviceColor(area.actionService)]),
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
                            Text(name.isEmpty ? area.name : name)
                                .font(.system(size: 20, weight: .semibold))
                                .foregroundColor(.white)
                                .shadow(color: .black.opacity(0.6), radius: 6, x: 0, y: 2)
                            Text("\(area.triggerService) → \(area.actionService)")
                                .font(.system(size: 14))
                                .foregroundColor(.white.opacity(0.9))
                                .shadow(color: .black.opacity(0.5), radius: 6, x: 0, y: 2)
                        }
                        Spacer()
                        ZStack {
                            Circle()
                                .fill(Color.white.opacity(0.12))
                                .frame(width: 48, height: 48)
                            Image(systemName: serviceIcon(area.actionService))
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

#Preview {
    EditAreaView(
        area: Area(
            id: "1",
            name: "Google Calendar → Gmail",
            description: "This AREA sends you an email via Gmail when a new event from your Google Calendar starts.",
            triggerService: "Google Calendar",
            actionService: "Gmail",
            isActive: true,
            isPublic: false,
            createdAt: "",
            updatedAt: "",
            userID: 1,
            triggerIconURL: nil,
            actionIconURL: nil,
            status: nil,
            triggerType: nil,
            actionType: nil,
            triggerConfig: nil,
            actionConfig: nil,
            conditions: nil,
            scheduleCron: nil,
            rateLimitPerMin: nil,
            dedupWindowSec: nil,
            retryMax: nil,
            retryBackoffMs: nil,
            lastRunStatus: nil,
            lastRunAt: nil,
            nextRunAt: nil,
            runCount: nil,
            lastError: nil,
            dedupKeyTemplate: nil,
            user: nil
        )
    )
}