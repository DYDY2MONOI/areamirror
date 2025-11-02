//
//  GenericConfigView.swift
//  area
//
//  Created to provide generic configuration UI for any service
//

import SwiftUI

struct GenericConfigView: View {
    let fields: [ServiceConfigField]
    @Binding var configValues: [String: ConfigFieldValue]
    let serviceName: String
    
    var body: some View {
        if !fields.isEmpty {
            VStack(alignment: .leading, spacing: 12) {
                ForEach(visibleFields) { field in
                    ConfigFieldView(
                        field: field,
                        value: Binding(
                            get: { configValues[field.key] ?? field.defaultValue ?? .string("") },
                            set: { configValues[field.key] = $0 }
                        )
                    )
                }
            }
        }
    }
    
    // Filter fields based on conditional logic (e.g., Twitter action mode)
    var visibleFields: [ServiceConfigField] {
        let normalizedServiceName = serviceName.lowercased().trimmingCharacters(in: .whitespaces)
        return fields.filter { field in
            // For Date Timer, hide eventDate and eventTime (handled by DatePicker)
            if normalizedServiceName == "date timer" {
                if field.key == "eventDate" || field.key == "eventTime" {
                    return false
                }
            }
            
            // For Twitter action, show conditional fields based on actionMode
            if normalizedServiceName == "twitter" {
                if let actionMode = configValues["actionMode"]?.stringValue {
                    if field.key == "tweetText" || field.key == "replyToTweetId" {
                        return actionMode == "tweet"
                    }
                    if field.key == "tweetId" {
                        return actionMode == "retweet"
                    }
                }
            }
            // For Telegram trigger, show conditional fields based on triggerType
            if normalizedServiceName == "telegram" {
                if let triggerType = configValues["triggerType"]?.stringValue {
                    if field.key == "keyword" {
                        return triggerType == "keyword_match"
                    }
                    if field.key == "command" {
                        return triggerType == "command_received"
                    }
                }
            }
            return true
        }
    }
}

struct ConfigFieldView: View {
    let field: ServiceConfigField
    @Binding var value: ConfigFieldValue
    
    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(field.label)
                .font(.system(size: 14, weight: .medium))
                .foregroundColor(.white)
            
            switch field.fieldType {
            case .text, .email:
                TextField(field.placeholder ?? "", text: Binding(
                    get: { value.stringValue },
                    set: { value = .string($0) }
                ))
                .textFieldStyle(CustomTextFieldStyle())
                .autocapitalization(.none)
                .disableAutocorrection(true)
                .keyboardType(field.fieldType == .email ? .emailAddress : .default)
                
            case .multiline:
                ZStack(alignment: .topLeading) {
                    RoundedRectangle(cornerRadius: 12)
                        .fill(AppColors.darkBackground)
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                        )
                    
                    if value.stringValue.isEmpty {
                        Text(field.placeholder ?? "")
                            .foregroundColor(.gray.opacity(0.6))
                            .padding(.horizontal, 12)
                            .padding(.vertical, 16)
                            .allowsHitTesting(false)
                    }
                    
                    TextEditor(text: Binding(
                        get: { value.stringValue },
                        set: { value = .string($0) }
                    ))
                    .foregroundColor(.white)
                    .frame(minHeight: 100)
                    .padding(4)
                    .scrollContentBackground(.hidden)
                    .background(Color.clear)
                }
                
            case .number:
                TextField(field.placeholder ?? "", text: Binding(
                    get: { value.stringValue },
                    set: { value = .string($0) }
                ))
                .textFieldStyle(CustomTextFieldStyle())
                .keyboardType(.decimalPad)
                
            case .select(let options):
                Picker(field.label, selection: Binding(
                    get: { value.stringValue },
                    set: { value = .string($0) }
                )) {
                    ForEach(options) { option in
                        Text(option.label).tag(option.value)
                    }
                }
                .pickerStyle(.menu)
                .foregroundColor(.white)
                .padding(.horizontal, 16)
                .padding(.vertical, 12)
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(AppColors.darkBackground)
                        .overlay(
                            RoundedRectangle(cornerRadius: 12)
                                .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                        )
                )
                
            case .toggle:
                Toggle("", isOn: Binding(
                    get: { value.boolValue },
                    set: { value = .bool($0) }
                ))
                .toggleStyle(SwitchToggleStyle(tint: AppColors.primaryBlue))
            }
            
            if let helperText = field.helperText {
                Text(helperText)
                    .font(.caption)
                    .foregroundColor(.gray)
                    .padding(.horizontal, 4)
            }
        }
    }
}

// Helper to convert ConfigFieldValue dictionary to [String: AnyCodable]
func configValuesToAnyCodable(_ values: [String: ConfigFieldValue]) -> [String: AnyCodable] {
    var result: [String: AnyCodable] = [:]
    for (key, value) in values {
        switch value {
        case .string(let str):
            result[key] = AnyCodable(str)
        case .bool(let bool):
            result[key] = AnyCodable(bool)
        }
    }
    return result
}

// Helper to convert [String: AnyCodable] to ConfigFieldValue dictionary
func anyCodableToConfigValues(_ values: [String: AnyCodable]?) -> [String: ConfigFieldValue] {
    guard let values = values else { return [:] }
    var result: [String: ConfigFieldValue] = [:]
    for (key, anyCodable) in values {
        if let stringValue = anyCodable.value as? String {
            result[key] = .string(stringValue)
        } else if let boolValue = anyCodable.value as? Bool {
            result[key] = .bool(boolValue)
        } else if let intValue = anyCodable.value as? Int {
            result[key] = .string(String(intValue))
        } else if let doubleValue = anyCodable.value as? Double {
            result[key] = .string(String(doubleValue))
        }
    }
    return result
}

