//
//  AppConfig.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation

struct AppConfig {
    
    static let baseURL = EnvironmentConfig.baseURL
    
    static let databaseTimeout: TimeInterval = 30.0
    
    static let appName = "AREA"
    static let appVersion = "1.0.0"
    
    static let requestTimeout: TimeInterval = 30.0
    static let maxRetryAttempts = 3
    
    static func getAPIEndpoint(_ path: String) -> String {
        return "\(baseURL)\(path)"
    }
    
    static func getWebSocketURL() -> String {
        return baseURL.replacingOccurrences(of: "http", with: "ws")
    }
}
