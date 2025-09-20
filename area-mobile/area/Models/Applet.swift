//
//  Applet.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct Applet: Identifiable {
    let id = UUID()
    let title: String
    let subtitle: String
    let description: String
    let icon: String
    let gradient: LinearGradient
    let type: AppletType
    let action: () -> Void
    
    enum AppletType {
        case create
        case blend
        case podcast
        case music
    }
}

extension Applet {
    static let sampleApplets: [Applet] = [
        Applet(
            title: "Gmail → Discord",
            subtitle: "Auto notification",
            description: "Send Discord message when you receive important emails",
            icon: "envelope.fill",
            gradient: LinearGradient(
                colors: [Color.red.opacity(0.8), Color.orange.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Gmail → Discord") }
        ),
        Applet(
            title: "Spotify → Twitter",
            subtitle: "Auto sharing",
            description: "Automatically tweet your favorite tracks",
            icon: "music.note",
            gradient: LinearGradient(
                colors: [Color.green.opacity(0.8), Color.blue.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Spotify → Twitter") }
        ),
        Applet(
            title: "GitHub → Slack",
            subtitle: "Dev notifications",
            description: "Notify team on Slack about new commits",
            icon: "hammer.fill",
            gradient: LinearGradient(
                colors: [Color.purple.opacity(0.8), Color.pink.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("GitHub → Slack") }
        ),
        
        Applet(
            title: "Weather → Telegram",
            subtitle: "Daily reminder",
            description: "Get weather forecast every morning on Telegram",
            icon: "cloud.sun.fill",
            gradient: LinearGradient(
                colors: [Color.blue.opacity(0.8), Color.cyan.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Weather → Telegram") }
        ),
        Applet(
            title: "Instagram → Dropbox",
            subtitle: "Auto backup",
            description: "Automatically backup your Instagram stories",
            icon: "camera.fill",
            gradient: LinearGradient(
                colors: [Color.pink.opacity(0.8), Color.purple.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Instagram → Dropbox") }
        ),
        Applet(
            title: "YouTube → Notion",
            subtitle: "Auto documentation",
            description: "Automatically add new videos to your knowledge base",
            icon: "play.rectangle.fill",
            gradient: LinearGradient(
                colors: [Color.red.opacity(0.8), Color.orange.opacity(0.6)],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("YouTube → Notion") }
        )
    ]
}
