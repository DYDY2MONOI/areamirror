//
//  AppletTest.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

extension Applet {
    static let testApplets: [Applet] = [
        Applet(
            title: "Test Gmail → Discord",
            subtitle: "Test notification",
            description: "Test description for debugging",
            icon: "envelope.fill",
            gradient: LinearGradient(
                colors: [Color.red, Color.orange],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Test Gmail → Discord") }
        ),
        Applet(
            title: "Test Spotify → Twitter",
            subtitle: "Test sharing",
            description: "Test description for debugging",
            icon: "music.note",
            gradient: LinearGradient(
                colors: [Color.green, Color.blue],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Test Spotify → Twitter") }
        ),
        Applet(
            title: "Test GitHub → Slack",
            subtitle: "Test notifications",
            description: "Test description for debugging",
            icon: "hammer.fill",
            gradient: LinearGradient(
                colors: [Color.purple, Color.pink],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            ),
            type: .create,
            action: { print("Test GitHub → Slack") }
        )
    ]
}
