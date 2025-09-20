//
//  Service.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct Service: Identifiable, Hashable {
    let id = UUID()
    let name: String
    let icon: String
    let color: Color
    let category: ServiceCategory
    
    enum ServiceCategory {
        case social
        case productivity
        case communication
        case entertainment
        case development
        case other
    }
}

extension Service {
    static let availableServices: [Service] = [
        Service(name: "Twitter", icon: "twitter", color: .blue, category: .social),
        Service(name: "Instagram", icon: "instagram", color: .pink, category: .social),
        Service(name: "Facebook", icon: "facebook", color: .blue, category: .social),
        Service(name: "LinkedIn", icon: "linkedin-brands-solid-full", color: .blue, category: .social),
        Service(name: "TikTok", icon: "tiktok", color: .white, category: .social),
        
        Service(name: "Discord", icon: "discord", color: .purple, category: .communication),
        Service(name: "Slack", icon: "slack", color: .green, category: .communication),
        Service(name: "Telegram", icon: "telegram", color: .blue, category: .communication),
        Service(name: "WhatsApp", icon: "whatsapp", color: .green, category: .communication),
        
        Service(name: "Gmail", icon: "gmail", color: .red, category: .productivity),
        Service(name: "Google Calendar", icon: "google-calendar", color: .blue, category: .productivity),
        Service(name: "Notion", icon: "notion", color: .gray, category: .productivity),
        Service(name: "Trello", icon: "trello", color: .blue, category: .productivity),
        Service(name: "Asana", icon: "asana", color: .orange, category: .productivity),
        Service(name: "Dropbox", icon: "dropbox", color: .blue, category: .productivity),
        
        Service(name: "Spotify", icon: "spotify", color: .green, category: .entertainment),
        Service(name: "YouTube", icon: "youtube", color: .red, category: .entertainment),
        Service(name: "Netflix", icon: "netflix", color: .red, category: .entertainment),
        Service(name: "Twitch", icon: "twitch", color: .purple, category: .entertainment),
        
        Service(name: "GitHub", icon: "github", color: .white, category: .development),
        Service(name: "GitLab", icon: "gitlab", color: .orange, category: .development),
        Service(name: "Jira", icon: "jira", color: .blue, category: .development),
        
        Service(name: "Weather", icon: "weather", color: .cyan, category: .other),
        Service(name: "News", icon: "news", color: .orange, category: .other),
        Service(name: "Reddit", icon: "reddit", color: .orange, category: .other)
    ]
    
    static func services(for category: ServiceCategory) -> [Service] {
        return availableServices.filter { $0.category == category }
    }
}
