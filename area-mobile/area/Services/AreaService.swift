//
//  AreaService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import SwiftUI

class AreaService: ObservableObject {
    static let shared = AreaService()
    
    @Published var userAreas: [Area] = []
    @Published var popularAreas: [AreaTemplate] = []
    @Published var recommendedAreas: [AreaTemplate] = []
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private init() {}
    
    func fetchUserAreas() async {
        guard let token = AuthService.shared.getAuthToken() else {
            DispatchQueue.main.async {
                self.errorMessage = "No authentication token"
            }
            return
        }
        
        DispatchQueue.main.async {
            self.isLoading = true
            self.errorMessage = nil
        }
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/user/me/areas")) else {
            DispatchQueue.main.async {
                self.errorMessage = "Invalid URL"
                self.isLoading = false
            }
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            let (data, response) = try await URLSession.shared.data(for: request)
            
            if let httpResponse = response as? HTTPURLResponse {
                if httpResponse.statusCode == 200 {
                    let areaResponse = try JSONDecoder().decode(AreaResponse.self, from: data)
                    print("✅ Fetched \(areaResponse.data.count) user areas")
                    DispatchQueue.main.async {
                        self.userAreas = areaResponse.data
                        self.isLoading = false
                    }
                } else {
                    DispatchQueue.main.async {
                        self.errorMessage = "Failed to fetch user areas"
                        self.isLoading = false
                    }
                }
            }
        } catch {
            print("❌ Error fetching user areas: \(error.localizedDescription)")
            DispatchQueue.main.async {
                self.errorMessage = "Network error: \(error.localizedDescription)"
                self.isLoading = false
            }
        }
    }
    
    func fetchPopularAreas() async {
        DispatchQueue.main.async {
            self.isLoading = true
            self.errorMessage = nil
        }
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas/popular")) else {
            DispatchQueue.main.async {
                self.errorMessage = "Invalid URL"
                self.isLoading = false
            }
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            let (data, response) = try await URLSession.shared.data(for: request)
            
            if let httpResponse = response as? HTTPURLResponse {
                if httpResponse.statusCode == 200 {
                    let areaResponse = try JSONDecoder().decode(AreaTemplateResponse.self, from: data)
                    print("✅ Fetched \(areaResponse.data.count) popular areas")
                    DispatchQueue.main.async {
                        self.popularAreas = areaResponse.data
                        self.isLoading = false
                    }
                } else {
                    DispatchQueue.main.async {
                        self.errorMessage = "Failed to fetch popular areas"
                        self.isLoading = false
                    }
                }
            }
        } catch {
            DispatchQueue.main.async {
                self.errorMessage = "Network error: \(error.localizedDescription)"
                self.isLoading = false
            }
        }
    }
    
    func fetchRecommendedAreas() async {
        DispatchQueue.main.async {
            self.isLoading = true
            self.errorMessage = nil
        }
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas/recommended")) else {
            DispatchQueue.main.async {
                self.errorMessage = "Invalid URL"
                self.isLoading = false
            }
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            let (data, response) = try await URLSession.shared.data(for: request)
            
            if let httpResponse = response as? HTTPURLResponse {
                if httpResponse.statusCode == 200 {
                    let areaResponse = try JSONDecoder().decode(AreaTemplateResponse.self, from: data)
                    print("✅ Fetched \(areaResponse.data.count) recommended areas")
                    DispatchQueue.main.async {
                        self.recommendedAreas = areaResponse.data
                        self.isLoading = false
                    }
                } else {
                    DispatchQueue.main.async {
                        self.errorMessage = "Failed to fetch recommended areas"
                        self.isLoading = false
                    }
                }
            }
        } catch {
            DispatchQueue.main.async {
                self.errorMessage = "Network error: \(error.localizedDescription)"
                self.isLoading = false
            }
        }
    }
    
    func fetchAllAreas() async {
        await withTaskGroup(of: Void.self) { group in
            group.addTask { await self.fetchUserAreas() }
            group.addTask { await self.fetchPopularAreas() }
            group.addTask { await self.fetchRecommendedAreas() }
        }
    }
}

struct AreaResponse: Codable {
    let data: [Area]
}

struct AreaTemplateResponse: Codable {
    let data: [AreaTemplate]
}

struct Area: Identifiable, Codable {
    let id: String
    let name: String
    let description: String
    let triggerService: String
    let actionService: String
    let isActive: Bool
    let isPublic: Bool
    let createdAt: String
    let updatedAt: String
    let userID: Int
    let triggerIconURL: String?
    let actionIconURL: String?
    
    enum CodingKeys: String, CodingKey {
        case id, name, description, isActive, isPublic
        case triggerService = "trigger_service"
        case actionService = "action_service"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case userID = "user_id"
        case triggerIconURL = "trigger_icon_url"
        case actionIconURL = "action_icon_url"
    }
}

struct AreaTemplate: Identifiable, Codable {
    let id: String
    let title: String
    let subtitle: String
    let description: String
    let icon: String
    let gradientClass: String
    let triggerService: String
    let actionService: String
    let triggerIconURL: String?
    let actionIconURL: String?
    let isActive: Bool
    
    enum CodingKeys: String, CodingKey {
        case id, title, subtitle, description, icon, isActive
        case gradientClass = "gradientClass"
        case triggerService = "triggerService"
        case actionService = "actionService"
        case triggerIconURL = "triggerIconUrl"
        case actionIconURL = "actionIconUrl"
    }
}