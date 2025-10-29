//
//  EnvironmentConfig.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation

enum AppEnvironment {
    case development
    case staging
    case production
    
    var baseURL: String {
        switch self {
        case .development:
            return "http://10.68.251.68:8080"
        case .staging:
            return "https://staging-api.area.com"
        case .production:
            return "https://api.area.com"
        }
    }
    
    var apiVersion: String {
        switch self {
        case .development, .staging:
            return "v1"
        case .production:
            return "v1"
        }
    }
    
    var isDebugMode: Bool {
        switch self {
        case .development:
            return true
        case .staging, .production:
            return false
        }
    }
}

struct EnvironmentConfig {
    static let current: AppEnvironment = {
        #if DEBUG
        return .development
        #else
        return .production
        #endif
    }()
    
    static var baseURL: String {
        return current.baseURL
    }
    
    static var apiVersion: String {
        return current.apiVersion
    }
    
    static var isDebugMode: Bool {
        return current.isDebugMode
    }
}
