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
        Service(name: "Twitter", icon: "bird.fill", color: .blue, category: .social),
        Service(name: "Instagram", icon: "camera.fill", color: .pink, category: .social),
        Service(name: "Facebook", icon: "f.circle.fill", color: .blue, category: .social),
        Service(name: "LinkedIn", icon: "briefcase.fill", color: .blue, category: .social),
        Service(name: "TikTok", icon: "music.note", color: .black, category: .social),
        
        Service(name: "Discord", icon: "message.fill", color: .purple, category: .communication),
        Service(name: "Slack", icon: "bubble.left.fill", color: .green, category: .communication),
        Service(name: "Telegram", icon: "paperplane.fill", color: .blue, category: .communication),
        Service(name: "WhatsApp", icon: "message.circle.fill", color: .green, category: .communication),
        
        Service(name: "Gmail", icon: "envelope.fill", color: .red, category: .productivity),
        Service(name: "Google Calendar", icon: "calendar", color: .blue, category: .productivity),
        Service(name: "Notion", icon: "doc.text.fill", color: .gray, category: .productivity),
        Service(name: "Trello", icon: "list.bullet.rectangle.fill", color: .blue, category: .productivity),
        Service(name: "Asana", icon: "checkmark.circle.fill", color: .orange, category: .productivity),
        Service(name: "Dropbox", icon: "folder.fill", color: .blue, category: .productivity),
        
        Service(name: "Spotify", icon: "music.note", color: .green, category: .entertainment),
        Service(name: "YouTube", icon: "play.rectangle.fill", color: .red, category: .entertainment),
        Service(name: "Netflix", icon: "tv.fill", color: .red, category: .entertainment),
        Service(name: "Twitch", icon: "gamecontroller.fill", color: .purple, category: .entertainment),
        
        Service(name: "GitHub", icon: "hammer.fill", color: .black, category: .development),
        Service(name: "GitLab", icon: "hammer.fill", color: .orange, category: .development),
        Service(name: "Jira", icon: "gear", color: .blue, category: .development),
        
        Service(name: "Weather", icon: "cloud.sun.fill", color: .cyan, category: .other),
        Service(name: "News", icon: "newspaper.fill", color: .orange, category: .other),
        Service(name: "Reddit", icon: "r.circle.fill", color: .orange, category: .other)
    ]
    
    static func services(for category: ServiceCategory) -> [Service] {
        return availableServices.filter { $0.category == category }
    }
}
