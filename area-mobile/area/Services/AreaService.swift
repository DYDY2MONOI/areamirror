//
//  AreaService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import SwiftUI

@MainActor
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
            self.errorMessage = "No authentication token"
            return
        }
        
        self.isLoading = true
        self.errorMessage = nil
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/user/me/areas")) else {
            self.errorMessage = "Invalid URL"
            self.isLoading = false
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
                    self.userAreas = areaResponse.data
                    self.isLoading = false
                } else {
                    self.errorMessage = "Failed to fetch user areas"
                    self.isLoading = false
                }
            }
        } catch {
            print("❌ Error fetching user areas: \(error.localizedDescription)")
            self.errorMessage = "Network error: \(error.localizedDescription)"
            self.isLoading = false
        }
    }
    
    func fetchPopularAreas() async {
        self.isLoading = true
        self.errorMessage = nil
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas/popular")) else {
            self.errorMessage = "Invalid URL"
            self.isLoading = false
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
                    self.popularAreas = areaResponse.data
                    self.isLoading = false
                } else {
                    self.errorMessage = "Failed to fetch popular areas"
                    self.isLoading = false
                }
            }
        } catch {
            self.errorMessage = "Network error: \(error.localizedDescription)"
            self.isLoading = false
        }
    }
    
    func fetchRecommendedAreas() async {
        self.isLoading = true
        self.errorMessage = nil
        
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas/recommended")) else {
            self.errorMessage = "Invalid URL"
            self.isLoading = false
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
                    self.recommendedAreas = areaResponse.data
                    self.isLoading = false
                } else {
                    self.errorMessage = "Failed to fetch recommended areas"
                    self.isLoading = false
                }
            }
        } catch {
            self.errorMessage = "Network error: \(error.localizedDescription)"
            self.isLoading = false
        }
    }
    
    func fetchAllAreas() async {
        await withTaskGroup(of: Void.self) { group in
            group.addTask { await self.fetchUserAreas() }
            group.addTask { await self.fetchPopularAreas() }
            group.addTask { await self.fetchRecommendedAreas() }
        }
    }

    // MARK: - Create/Update
    struct CreateOrUpdateAreaRequest: Codable {
        let name: String
        let description: String
        let triggerService: String
        let triggerType: String
        let actionService: String
        let actionType: String
        let triggerConfig: [String: AnyCodable]
        let actionConfig: [String: AnyCodable]
    }

    struct SingleAreaResponse: Codable {
        let data: Area
    }

    func createArea(payload: CreateOrUpdateAreaRequest) async throws -> Area {
        guard let token = AuthService.shared.getAuthToken() else {
            throw AreaServiceError.unauthorized
        }
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas")) else {
            throw AreaServiceError.invalidURL
        }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        let encoder = JSONEncoder()
        request.httpBody = try encoder.encode(payload)
        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200...299).contains(http.statusCode) else {
            let msg = String(data: data, encoding: .utf8) ?? "Server error"
            throw AreaServiceError.server(msg)
        }
        let decoded = try JSONDecoder().decode(SingleAreaResponse.self, from: data)
        self.userAreas.insert(decoded.data, at: 0)
        return decoded.data
    }

    func updateArea(areaId: String, payload: CreateOrUpdateAreaRequest) async throws -> Area {
        guard let token = AuthService.shared.getAuthToken() else {
            throw AreaServiceError.unauthorized
        }
        guard let url = URL(string: AppConfig.getAPIEndpoint("/areas/\(areaId)")) else {
            throw AreaServiceError.invalidURL
        }
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        let encoder = JSONEncoder()
        request.httpBody = try encoder.encode(payload)
        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200...299).contains(http.statusCode) else {
            let msg = String(data: data, encoding: .utf8) ?? "Server error"
            throw AreaServiceError.server(msg)
        }
        let decoded = try JSONDecoder().decode(SingleAreaResponse.self, from: data)
        if let idx = self.userAreas.firstIndex(where: { $0.id == decoded.data.id }) {
            self.userAreas[idx] = decoded.data
        }
        return decoded.data
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

// MARK: - Helpers
enum AreaServiceError: Error, LocalizedError {
    case unauthorized
    case invalidURL
    case server(String)
    
    var errorDescription: String? {
        switch self {
        case .unauthorized: return "No authentication token"
        case .invalidURL: return "Invalid URL"
        case .server(let msg): return msg
        }
    }
}

// A simple type-erased wrapper to encode heterogeneous values in dictionaries
struct AnyCodable: Codable {
    let value: Any
    init(_ value: Any) { self.value = value }
    
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        if let intVal = try? container.decode(Int.self) { value = intVal; return }
        if let doubleVal = try? container.decode(Double.self) { value = doubleVal; return }
        if let boolVal = try? container.decode(Bool.self) { value = boolVal; return }
        if let stringVal = try? container.decode(String.self) { value = stringVal; return }
        if let dictVal = try? container.decode([String: AnyCodable].self) { value = dictVal.mapValues { $0.value }; return }
        if let arrayVal = try? container.decode([AnyCodable].self) { value = arrayVal.map { $0.value }; return }
        throw DecodingError.dataCorruptedError(in: container, debugDescription: "Unsupported type")
    }
    
    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        switch value {
        case let intVal as Int: try container.encode(intVal)
        case let doubleVal as Double: try container.encode(doubleVal)
        case let boolVal as Bool: try container.encode(boolVal)
        case let stringVal as String: try container.encode(stringVal)
        case let dictVal as [String: Any]: try container.encode(dictVal.mapValues { AnyCodable($0) })
        case let arrayVal as [Any]: try container.encode(arrayVal.map { AnyCodable($0) })
        default:
            let ctx = EncodingError.Context(codingPath: container.codingPath, debugDescription: "Unsupported type")
            throw EncodingError.invalidValue(value, ctx)
        }
    }
}
