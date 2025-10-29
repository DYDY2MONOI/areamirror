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
        Service(name: "GitHub", icon: "github", color: Color(hex: "#24292e"), category: .development),
        Service(name: "Google", icon: "gmail", color: Color(hex: "#4285f4"), category: .productivity),
        Service(name: "Facebook", icon: "facebook", color: Color(hex: "#1877f2"), category: .social),
        Service(name: "Spotify", icon: "spotify", color: Color(hex: "#1db954"), category: .entertainment),
        Service(name: "Twitter / X", icon: "twitter", color: Color(hex: "#1DA1F2"), category: .social)
    ]
    
    static func services(for category: ServiceCategory) -> [Service] {
        return availableServices.filter { $0.category == category }
    }
}
