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
    
    @State private var triggerConfigValues: [String: ConfigFieldValue] = [:]
    @State private var actionConfigValues: [String: ConfigFieldValue] = [:]
    
    @State private var eventDateTime: Date = Date()
    
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
            print("📱 Trigger service name: '\(area.triggerService)'")
            triggerConfigValues = anyCodableToConfigValues(triggerConfig)
            
            let normalizedTriggerService = area.triggerService.lowercased().trimmingCharacters(in: .whitespaces)
            if normalizedTriggerService == "date timer" {
                let eventDate = triggerConfigValues["eventDate"]?.stringValue ?? ""
                let eventTime = triggerConfigValues["eventTime"]?.stringValue ?? ""
                if let parsed = Self.parseCalendarDateTime(dateString: eventDate, timeString: eventTime) {
                    _eventDateTime = State(initialValue: parsed)
                }
            }
        } else {
            let triggerFields = ServiceConfigMetadata.triggerConfigFields(for: area.triggerService)
            for field in triggerFields {
                if let defaultValue = field.defaultValue {
                    triggerConfigValues[field.key] = defaultValue
                }
            }
        }
        
        if let actionConfig = area.actionConfig {
            print("📱 Loading action config for existing area: \(area.name)")
            print("📱 Action service name: '\(area.actionService)'")
            actionConfigValues = anyCodableToConfigValues(actionConfig)
        } else {
            let actionFields = ServiceConfigMetadata.actionConfigFields(for: area.actionService)
            for field in actionFields {
                if let defaultValue = field.defaultValue {
                    actionConfigValues[field.key] = defaultValue
                }
            }
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
                        
                        #if DEBUG
                        VStack(alignment: .leading, spacing: 4) {
                            Text("DEBUG: triggerService = '\(area.triggerService)'")
                                .foregroundColor(.green)
                                .font(.caption)
                            Text("DEBUG: actionService = '\(area.actionService)'")
                                .foregroundColor(.green)
                                .font(.caption)
                        }
                        .padding(.horizontal, 20)
                        #endif
                        
                        if let triggerFields = getTriggerConfigFields(for: area.triggerService), !triggerFields.isEmpty {
                            genericTriggerConfigCard(serviceName: area.triggerService, fields: triggerFields)
                        }
                        
                        if let actionFields = getActionConfigFields(for: area.actionService), !actionFields.isEmpty {
                            genericActionConfigCard(serviceName: area.actionService, fields: actionFields)
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
    
    @ViewBuilder
    private func genericTriggerConfigCard(serviceName: String, fields: [ServiceConfigField]) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: serviceIcon(serviceName))
                    .foregroundColor(.white)
                Text("\(serviceName) Trigger")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(.white)
            }
            
            if serviceName.lowercased() == "date timer" {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Date & time")
                        .foregroundColor(.white)
                        .font(.system(size: 14, weight: .medium))
                    DatePicker("", selection: $eventDateTime, displayedComponents: [.date, .hourAndMinute])
                        .labelsHidden()
                        .datePickerStyle(.compact)
                        .onChange(of: eventDateTime) { oldValue, newValue in
                            let (dateStr, timeStr) = Self.formatCalendarStrings(from: newValue)
                            triggerConfigValues["eventDate"] = .string(dateStr)
                            triggerConfigValues["eventTime"] = .string(timeStr)
                        }
                }
            }
            
            GenericConfigView(
                fields: fields,
                configValues: $triggerConfigValues,
                serviceName: serviceName
            )
        }
        .cardStyle()
        .padding(.horizontal, 20)
    }
    
    @ViewBuilder
    private func genericActionConfigCard(serviceName: String, fields: [ServiceConfigField]) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            HStack(spacing: 8) {
                Image(systemName: serviceIcon(serviceName))
                    .foregroundColor(.white)
                Text("\(serviceName) Action")
                    .font(.system(size: 18, weight: .semibold))
                    .foregroundColor(.white)
            }
            
            GenericConfigView(
                fields: fields,
                configValues: $actionConfigValues,
                serviceName: serviceName
            )
        }
        .cardStyle()
        .padding(.horizontal, 20)
    }
    
    private func getTriggerConfigFields(for serviceName: String) -> [ServiceConfigField]? {
        let fields = ServiceConfigMetadata.triggerConfigFields(for: serviceName)
        return fields.isEmpty ? nil : fields
    }
    
    private func getActionConfigFields(for serviceName: String) -> [ServiceConfigField]? {
        let fields = ServiceConfigMetadata.actionConfigFields(for: serviceName)
        return fields.isEmpty ? nil : fields
    }
    
    private var canSave: Bool {
        if name.isEmpty { return false }
        
        if let triggerFields = getTriggerConfigFields(for: area.triggerService) {
            for field in triggerFields where field.isRequired {
                let value = triggerConfigValues[field.key]?.stringValue ?? ""
                if value.isEmpty {
                    return false
                }
            }
        }
        
        if let actionFields = getActionConfigFields(for: area.actionService) {
            for field in actionFields where field.isRequired {
                let value = actionConfigValues[field.key]?.stringValue ?? ""
                if value.isEmpty {
                    return false
                }
            }
            
            let normalizedActionService = area.actionService.lowercased().trimmingCharacters(in: .whitespaces)
            if normalizedActionService == "twitter" {
                if actionConfigValues["actionMode"]?.stringValue == "tweet" {
                    let tweetText = actionConfigValues["tweetText"]?.stringValue ?? ""
                    if tweetText.isEmpty {
                        return false
                    }
                } else if actionConfigValues["actionMode"]?.stringValue == "retweet" {
                    let tweetId = actionConfigValues["tweetId"]?.stringValue ?? ""
                    if tweetId.isEmpty {
                        return false
                    }
                }
            }
            
            if normalizedActionService == "spotify" {
                let playlistId = actionConfigValues["playlistId"]?.stringValue ?? ""
                let range = actionConfigValues["range"]?.stringValue ?? ""
                let urlColumn = actionConfigValues["urlColumn"]?.stringValue ?? ""
                if playlistId.isEmpty || range.isEmpty || urlColumn.isEmpty {
                    return false
                }
            }
        }
        
        return true
    }
    
    private func saveArea() {
        isSaving = true
        errorMessage = nil
        
        Task {
            do {
                var triggerConfig: [String: AnyCodable] = [:]
                let normalizedTriggerService = area.triggerService.lowercased().trimmingCharacters(in: .whitespaces)
                
                if normalizedTriggerService == "date timer" {
                    let (dateStr, timeStr) = Self.formatCalendarStrings(from: eventDateTime)
                    triggerConfigValues["eventDate"] = .string(dateStr)
                    triggerConfigValues["eventTime"] = .string(timeStr)
                }
                
                triggerConfig = configValuesToAnyCodable(triggerConfigValues)
                
                var actionConfig: [String: AnyCodable] = [:]
                actionConfig = configValuesToAnyCodable(actionConfigValues)
                
                let fullPayload = AreaService.CreateOrUpdateAreaRequest(
                    name: name,
                    description: description,
                    triggerService: area.triggerService,
                    triggerType: area.triggerService == "Date Timer" ? "Event" : "Webhook",
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
        case "Twitter", "Twitter / X": return "PostTweet"
        default: return "Action"
        }
    }
}

extension EditAreaView {
    private static func parseCalendarDateTime(dateString: String, timeString: String) -> Date? {
        if !timeString.isEmpty {
            let iso = ISO8601DateFormatter()
            if let d = iso.date(from: timeString) {
                return d
            }
        }

        if !dateString.isEmpty {
            let timeRegex = try! NSRegularExpression(pattern: "T(\\d{2}:\\d{2})", options: [])
            var hhmm: String?
            if let match = timeRegex.firstMatch(in: timeString, options: [], range: NSRange(location: 0, length: (timeString as NSString).length)) {
                let r = match.range(at: 1)
                if r.location != NSNotFound, let rr = Range(r, in: timeString) {
                    hhmm = String(timeString[rr])
                }
            }
            if hhmm == nil {
                let simple = try! NSRegularExpression(pattern: "^\\d{2}:\\d{2}$", options: [])
                if simple.firstMatch(in: timeString, options: [], range: NSRange(location: 0, length: (timeString as NSString).length)) != nil {
                    hhmm = timeString
                }
            }

            let parts = dateString.split(separator: "-")
            guard parts.count == 3, let year = Int(parts[0]), let month = Int(parts[1]), let day = Int(parts[2]) else {
                return nil
            }
            var comps = DateComponents()
            comps.year = year
            comps.month = month
            comps.day = day
            if let hhmm = hhmm {
                let t = hhmm.split(separator: ":")
                if t.count >= 2, let h = Int(t[0]), let m = Int(t[1]) {
                    comps.hour = h
                    comps.minute = m
                }
            }
            comps.second = 0
            comps.timeZone = TimeZone.current
            return Calendar.current.date(from: comps)
        }

        return nil
    }

    private static func formatCalendarStrings(from date: Date) -> (String, String) {
        let df = DateFormatter()
        df.locale = Locale(identifier: "en_US_POSIX")
        df.timeZone = TimeZone.current
        df.dateFormat = "yyyy-MM-dd"
        let dateStr = df.string(from: date)

        let iso = ISO8601DateFormatter()
        iso.timeZone = TimeZone.current
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
            name: "Date Timer → Gmail",
            description: "This AREA sends you an email via Gmail when a new event from your Date Timer starts.",
            triggerService: "Date Timer",
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