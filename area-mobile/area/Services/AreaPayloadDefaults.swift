//
//  AreaPayloadDefaults.swift
//  area
//
//  Created by Codex on 03/12/2025.
//

import Foundation

enum AreaPayloadDefaults {
    static func triggerType(for service: String) -> String {
        switch service {
        case "Google Calendar":
            return "Event"
        case "Google Sheets":
            return "SpreadsheetChange"
        case "Spotify":
            return "Playback"
        case "Telegram":
            return "message_received"
        default:
            return "Webhook"
        }
    }

    static func actionType(for service: String) -> String {
        switch service {
        case "Gmail":
            return "SendEmail"
        default:
            return "Action"
        }
    }

    static func triggerConfig(for service: String, actionName: String) -> [String: AnyCodable] {
        switch service {
        case "Google Calendar":
            return [
                "eventDate": AnyCodable(""),
                "eventTime": AnyCodable(""),
                "eventTitle": AnyCodable(actionName.isEmpty ? "Event" : actionName),
                "calendarId": AnyCodable("primary")
            ]
        case "Google Sheets":
            return [
                "spreadsheetId": AnyCodable(""),
                "sheetName": AnyCodable(""),
                "range": AnyCodable("Sheet1!A1:D"),
                "hasHeader": AnyCodable(true)
            ]
        case "Weather":
            return [
                "city": AnyCodable(""),
                "temperature": AnyCodable(0),
                "condition": AnyCodable("")
            ]
        case "Telegram":
            return [
                "chatId": AnyCodable(""),
                "triggerType": AnyCodable("message_received")
            ]
        default:
            return [:]
        }
    }

    static func actionConfig(for service: String) -> [String: AnyCodable] {
        switch service {
        case "Gmail":
            return [
                "toEmail": AnyCodable(""),
                "subject": AnyCodable(""),
                "body": AnyCodable("")
            ]
        case "Slack":
            return [
                "channel": AnyCodable(""),
                "message": AnyCodable("")
            ]
        case "Discord":
            return [
                "channel": AnyCodable(""),
                "message": AnyCodable("")
            ]
        default:
            return [:]
        }
    }
}
