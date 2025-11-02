//
//  ServiceConfigMetadata.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import Foundation

enum ServiceConfigFieldType: Equatable {
    case text
    case email
    case multiline
    case number
    case select(options: [ServiceConfigFieldOption])
    case toggle
}

struct ServiceConfigFieldOption: Identifiable, Equatable {
    let value: String
    let label: String
    
    var id: String { value }
}

enum ConfigFieldValue: Equatable {
    case string(String)
    case bool(Bool)
    
    var stringValue: String {
        switch self {
        case .string(let value):
            return value
        case .bool(let value):
            return value ? "true" : "false"
        }
    }
    
    var boolValue: Bool {
        switch self {
        case .string(let value):
            return (value as NSString).boolValue
        case .bool(let value):
            return value
        }
    }
    
    static func from(any value: Any) -> ConfigFieldValue? {
        switch value {
        case let bool as Bool:
            return .bool(bool)
        case let int as Int:
            return .string(String(int))
        case let double as Double:
            var formatted = String(double)
            if formatted.hasSuffix(".0") {
                formatted.removeLast(2)
            }
            return .string(formatted)
        case let string as String:
            return .string(string)
        default:
            return nil
        }
    }
}

struct ServiceConfigField: Identifiable, Equatable {
    let key: String
    let label: String
    let placeholder: String?
    let helperText: String?
    let fieldType: ServiceConfigFieldType
    let isRequired: Bool
    let defaultValue: ConfigFieldValue?
    
    var id: String { key }
}

struct ServiceConfigMetadata {
    static func triggerConfigFields(for serviceName: String) -> [ServiceConfigField] {
        let normalized = serviceName.lowercased().trimmingCharacters(in: .whitespaces)
        switch normalized {
        case "google calendar", "date timer":
            return [
                ServiceConfigField(
                    key: "eventDate",
                    label: "Event Date",
                    placeholder: "2024-12-25",
                    helperText: "Date for the event (YYYY-MM-DD)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "eventTime",
                    label: "Event Time",
                    placeholder: "14:00",
                    helperText: "Time for the event (HH:MM)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "eventTitle",
                    label: "Event Title",
                    placeholder: "Meeting with team",
                    helperText: "Title of the calendar event",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "calendarId",
                    label: "Calendar ID",
                    placeholder: "primary",
                    helperText: "Use 'primary' for your main calendar",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("primary")
                )
            ]
        case "google sheets":
            return [
                ServiceConfigField(
                    key: "spreadsheetId",
                    label: "Spreadsheet ID",
                    placeholder: "1A2B3C4D...",
                    helperText: "Copy the ID from the Sheets URL (between /d/ and /edit)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "sheetName",
                    label: "Sheet Name",
                    placeholder: "Sheet1",
                    helperText: "Name of the sheet to monitor",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "range",
                    label: "Range (A1 notation)",
                    placeholder: "Sheet1!A1:D",
                    helperText: "Range to monitor (format: Sheet1!A1:D)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("Sheet1!A1:D")
                ),
                ServiceConfigField(
                    key: "hasHeader",
                    label: "Has Header Row",
                    placeholder: nil,
                    helperText: "Check if the first row contains headers",
                    fieldType: .toggle,
                    isRequired: false,
                    defaultValue: .bool(true)
                )
            ]
        case "weather":
            return [
                ServiceConfigField(
                    key: "city",
                    label: "City",
                    placeholder: "Paris",
                    helperText: "City to monitor weather for",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "temperature",
                    label: "Temperature Threshold",
                    placeholder: "25",
                    helperText: "Temperature in Celsius",
                    fieldType: .number,
                    isRequired: false,
                    defaultValue: .string("0")
                ),
                ServiceConfigField(
                    key: "condition",
                    label: "Weather Condition",
                    placeholder: "Clear",
                    helperText: "Specific weather condition (optional)",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                )
            ]
        case "telegram":
            return [
                ServiceConfigField(
                    key: "chatId",
                    label: "Chat ID",
                    placeholder: "123456789",
                    helperText: "Your Telegram chat ID. Get it from @userinfobot",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "triggerType",
                    label: "Trigger Type",
                    placeholder: nil,
                    helperText: "When should this trigger?",
                    fieldType: .select(options: [
                        ServiceConfigFieldOption(value: "message_received", label: "Any Message"),
                        ServiceConfigFieldOption(value: "keyword_match", label: "Keyword Match"),
                        ServiceConfigFieldOption(value: "command_received", label: "Command Received")
                    ]),
                    isRequired: true,
                    defaultValue: .string("message_received")
                ),
                ServiceConfigField(
                    key: "keyword",
                    label: "Keyword",
                    placeholder: "urgent",
                    helperText: "Messages containing this keyword will trigger the area",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "command",
                    label: "Command",
                    placeholder: "/start",
                    helperText: "This specific command will trigger the area (e.g., /start, /help)",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                )
            ]
        case "twitter", "twitter / x":
            return [
                ServiceConfigField(
                    key: "monitorType",
                    label: "Monitor Type",
                    placeholder: nil,
                    helperText: "What to monitor on Twitter",
                    fieldType: .select(options: [
                        ServiceConfigFieldOption(value: "mentions", label: "Mentions"),
                        ServiceConfigFieldOption(value: "likes", label: "Likes"),
                        ServiceConfigFieldOption(value: "retweets", label: "Retweets"),
                        ServiceConfigFieldOption(value: "followers", label: "New Followers")
                    ]),
                    isRequired: true,
                    defaultValue: .string("mentions")
                ),
                ServiceConfigField(
                    key: "keyword",
                    label: "Keyword (optional)",
                    placeholder: "urgent",
                    helperText: "Filter mentions by keyword (optional)",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "includeRetweets",
                    label: "Include Retweets",
                    placeholder: nil,
                    helperText: "Include retweets in mentions",
                    fieldType: .toggle,
                    isRequired: false,
                    defaultValue: .bool(false)
                )
            ]
        case "github":
            return [
                ServiceConfigField(
                    key: "repositoryId",
                    label: "Repository ID",
                    placeholder: "123456789",
                    helperText: "GitHub repository ID to monitor",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "notificationTypes",
                    label: "Event Types",
                    placeholder: nil,
                    helperText: "Select GitHub events to monitor (multiple allowed)",
                    fieldType: .text, // Will be handled as comma-separated or array
                    isRequired: true,
                    defaultValue: .string("push")
                )
            ]
        case "google drive":
            return [
                ServiceConfigField(
                    key: "folderId",
                    label: "Folder ID",
                    placeholder: "1A2B3C4D...",
                    helperText: "Copy the ID from the Drive folder URL (between /folders/ and the end)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                )
            ]
        default:
            return []
        }
    }
    
    static func actionConfigFields(for serviceName: String) -> [ServiceConfigField] {
        let normalized = serviceName.lowercased().trimmingCharacters(in: .whitespaces)
        switch normalized {
        case "gmail":
            return [
                ServiceConfigField(
                    key: "toEmail",
                    label: "Send Email To",
                    placeholder: "your-email@gmail.com",
                    helperText: "Email address to send to",
                    fieldType: .email,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "subject",
                    label: "Email Subject",
                    placeholder: "Reminder: {{eventTitle}}",
                    helperText: "Use {{eventTitle}}, {{eventTime}}, {{areaName}} as variables",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "body",
                    label: "Email Body",
                    placeholder: "Hello! This is a reminder...",
                    helperText: "Email body content with template variables",
                    fieldType: .multiline,
                    isRequired: true,
                    defaultValue: .string("")
                )
            ]
        case "slack":
            return [
                ServiceConfigField(
                    key: "channel",
                    label: "Channel",
                    placeholder: "#general",
                    helperText: "Slack channel name or ID",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "message",
                    label: "Message",
                    placeholder: "Notification message",
                    helperText: "Message to send to Slack",
                    fieldType: .multiline,
                    isRequired: true,
                    defaultValue: .string("")
                )
            ]
        case "discord":
            return [
                ServiceConfigField(
                    key: "webhookUrl",
                    label: "Webhook URL",
                    placeholder: "https://discord.com/api/webhooks/...",
                    helperText: "Discord webhook URL (Server Settings → Integrations → Webhooks)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "message",
                    label: "Message",
                    placeholder: "Message to send to Discord",
                    helperText: "Message content with template variables",
                    fieldType: .multiline,
                    isRequired: true,
                    defaultValue: .string("")
                )
            ]
        case "twitter", "twitter / x":
            return [
                ServiceConfigField(
                    key: "actionMode",
                    label: "Action Type",
                    placeholder: nil,
                    helperText: "Post a tweet or retweet",
                    fieldType: .select(options: [
                        ServiceConfigFieldOption(value: "tweet", label: "Post Tweet"),
                        ServiceConfigFieldOption(value: "retweet", label: "Retweet")
                    ]),
                    isRequired: true,
                    defaultValue: .string("tweet")
                ),
                ServiceConfigField(
                    key: "tweetText",
                    label: "Tweet Text",
                    placeholder: "Thanks for the mention @{{tweetAuthorUsername}}! ",
                    helperText: "Tweet content (max 280 chars). Use {{tweetText}}, {{tweetAuthorUsername}}, etc.",
                    fieldType: .multiline,
                    isRequired: false,
                    defaultValue: .string("Thanks for the mention @{{tweetAuthorUsername}}! ")
                ),
                ServiceConfigField(
                    key: "replyToTweetId",
                    label: "Reply To Tweet ID (optional)",
                    placeholder: "{{tweetId}}",
                    helperText: "Tweet ID to reply to. Use {{tweetId}} for the triggering tweet",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "tweetId",
                    label: "Tweet ID to Retweet",
                    placeholder: "{{tweetId}}",
                    helperText: "Tweet ID to retweet. Use {{tweetId}} for the triggering tweet",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                )
            ]
        case "telegram":
            return [
                ServiceConfigField(
                    key: "chatId",
                    label: "Chat ID",
                    placeholder: "8481009224",
                    helperText: "Your Telegram chat ID. Get it from @userinfobot",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "message",
                    label: "Message",
                    placeholder: " Notification from {{areaName}}",
                    helperText: "Message to send. Use {{areaName}}, {{triggerTime}}, etc. Supports Markdown (*bold*, _italic_)",
                    fieldType: .multiline,
                    isRequired: true,
                    defaultValue: .string("")
                )
            ]
        case "spotify":
            return [
                ServiceConfigField(
                    key: "playlistId",
                    label: "Playlist ID",
                    placeholder: "37i9dQZF1DXcBWIGoYBM5M",
                    helperText: "The playlist ID is the final part of the Spotify URL (e.g., https://open.spotify.com/playlist/37i9dQZF...)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "spreadsheetId",
                    label: "Spreadsheet ID (optional)",
                    placeholder: "Leave empty to use trigger's spreadsheet",
                    helperText: "If left empty, the automation will use the ID configured in the Google Sheets trigger",
                    fieldType: .text,
                    isRequired: false,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "range",
                    label: "Sheet Range",
                    placeholder: "Sheet1!A2:C",
                    helperText: "Range to synchronize (recommended: start at row 2 to exclude header)",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("")
                ),
                ServiceConfigField(
                    key: "urlColumn",
                    label: "Spotify Link Column",
                    placeholder: "SpotifyLink",
                    helperText: "Name of the column containing Spotify links (URL or URI). Can be a letter (e.g., C) or a header name",
                    fieldType: .text,
                    isRequired: true,
                    defaultValue: .string("SpotifyLink")
                ),
                ServiceConfigField(
                    key: "hasHeader",
                    label: "Has Header Row",
                    placeholder: nil,
                    helperText: "Check if the first row contains headers",
                    fieldType: .toggle,
                    isRequired: false,
                    defaultValue: .bool(true)
                )
            ]
        default:
            return []
        }
    }
}
