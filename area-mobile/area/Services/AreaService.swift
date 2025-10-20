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
        var name: String?
        var description: String?
        var triggerService: String?
        var triggerType: String?
        var actionService: String?
        var actionType: String?
        var triggerConfig: [String: AnyCodable]?
        var actionConfig: [String: AnyCodable]?
        var isActive: Bool?
        
        enum CodingKeys: String, CodingKey {
            case name, description
            case triggerService = "triggerService"
            case triggerType = "triggerType"
            case actionService = "actionService"
            case actionType = "actionType"
            case triggerConfig = "triggerConfig"
            case actionConfig = "actionConfig"
            case isActive
        }
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
        let bodyData = try encoder.encode(payload)
        #if DEBUG
        if let json = String(data: bodyData, encoding: .utf8) { print("➡️ createArea body: \(json)") }
        #endif
        request.httpBody = bodyData
        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200...299).contains(http.statusCode) else {
            let msg = String(data: data, encoding: .utf8) ?? "Server error"
            throw AreaServiceError.server(msg)
        }
        do {
            let decoded = try JSONDecoder().decode(SingleAreaResponse.self, from: data)
            self.userAreas.insert(decoded.data, at: 0)
            return decoded.data
        } catch {
            if let lenient = try? Self.parseAreaLenient(from: data) {
                self.userAreas.insert(lenient, at: 0)
                return lenient
            }
            let body = String(data: data, encoding: .utf8) ?? "<no-body>"
            throw AreaServiceError.server("Decoding error: \(Self.describeDecodingError(error)) | Body: \(body)")
        }
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
        encoder.outputFormatting = .prettyPrinted
        
        let bodyData = try encoder.encode(payload)
        #if DEBUG
        if let json = String(data: bodyData, encoding: .utf8) { print("➡️ updateArea body: \(json)") }
        #endif
        request.httpBody = bodyData
        
        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200...299).contains(http.statusCode) else {
            let msg = String(data: data, encoding: .utf8) ?? "Server error"
            throw AreaServiceError.server(msg)
        }
        do {
            let decoded = try JSONDecoder().decode(SingleAreaResponse.self, from: data)
            if let idx = self.userAreas.firstIndex(where: { $0.id == decoded.data.id }) {
                self.userAreas[idx] = decoded.data
            }
            return decoded.data
        } catch {
            if let lenient = try? Self.parseAreaLenient(from: data) {
                if let idx = self.userAreas.firstIndex(where: { $0.id == lenient.id }) {
                    self.userAreas[idx] = lenient
                }
                return lenient
            }
            let body = String(data: data, encoding: .utf8) ?? "<no-body>"
            throw AreaServiceError.server("Decoding error: \(Self.describeDecodingError(error)) | Body: \(body)")
        }
    }
    
    // MARK: - Lenient decoding helpers
    private static func parseAreaLenient(from data: Data) throws -> Area? {
        guard
            let root = try JSONSerialization.jsonObject(with: data) as? [String: Any],
            let area = root["data"] as? [String: Any]
        else { return nil }
        
        // Create a clean dictionary with proper field mapping
        var cleanArea: [String: Any] = [:]
        
        // Map snake_case to camelCase for required fields
        cleanArea["id"] = area["id"] as? String ?? UUID().uuidString
        cleanArea["name"] = area["name"] as? String ?? "Untitled"
        cleanArea["description"] = area["description"] as? String ?? ""
        cleanArea["triggerService"] = area["trigger_service"] as? String ?? "Unknown"
        cleanArea["actionService"] = area["action_service"] as? String ?? "Unknown"
        cleanArea["isActive"] = area["is_active"] as? Bool ?? true
        cleanArea["isPublic"] = area["is_public"] as? Bool ?? true
        cleanArea["createdAt"] = area["created_at"] as? String ?? ""
        cleanArea["updatedAt"] = area["updated_at"] as? String ?? ""
        cleanArea["userID"] = area["user_id"] as? Int ?? 0
        
        // Map optional fields
        cleanArea["triggerIconURL"] = area["trigger_icon_url"]
        cleanArea["actionIconURL"] = area["action_icon_url"]
        cleanArea["status"] = area["status"]
        cleanArea["triggerType"] = area["trigger_type"]
        cleanArea["actionType"] = area["action_type"]
        cleanArea["triggerConfig"] = area["trigger_config"] ?? [:]
        cleanArea["actionConfig"] = area["action_config"] ?? [:]
        cleanArea["conditions"] = area["conditions"] ?? []
        cleanArea["scheduleCron"] = area["schedule_cron"]
        cleanArea["rateLimitPerMin"] = area["rate_limit_per_min"]
        cleanArea["dedupWindowSec"] = area["dedup_window_sec"]
        cleanArea["retryMax"] = area["retry_max"]
        cleanArea["retryBackoffMs"] = area["retry_backoff_ms"]
        cleanArea["lastRunStatus"] = area["last_run_status"]
        cleanArea["lastRunAt"] = area["last_run_at"]
        cleanArea["nextRunAt"] = area["next_run_at"]
        cleanArea["runCount"] = area["run_count"] ?? 0
        cleanArea["lastError"] = area["last_error"]
        cleanArea["dedupKeyTemplate"] = area["dedup_key_template"]
        
        // Handle user object
        if let user = area["user"] as? [String: Any] {
            var cleanUser: [String: Any] = [:]
            cleanUser["id"] = user["id"] as? Int ?? 0
            cleanUser["email"] = user["email"] as? String ?? ""
            cleanUser["firstName"] = user["first_name"]
            cleanUser["lastName"] = user["last_name"]
            cleanUser["createdAt"] = user["created_at"] as? String ?? ""
            cleanUser["updatedAt"] = user["updated_at"] as? String ?? ""
            cleanUser["phone"] = user["phone"]
            cleanUser["birthday"] = user["birthday"]
            cleanUser["gender"] = user["gender"]
            cleanUser["country"] = user["country"]
            cleanUser["lang"] = user["lang"]
            cleanUser["loginProvider"] = user["login_provider"]
            cleanUser["profileImage"] = user["profile_image"]
            cleanUser["githubID"] = user["github_id"]
            cleanUser["githubUsername"] = user["github_username"]
            cleanUser["googleID"] = user["google_id"]
            cleanUser["googleEmail"] = user["google_email"]
            cleanUser["facebookID"] = user["facebook_id"]
            cleanUser["facebookEmail"] = user["facebook_email"]
            cleanUser["role"] = user["role"]
            cleanUser["isActive"] = user["is_active"]
            cleanArea["user"] = cleanUser
        }
        
        let json = try JSONSerialization.data(withJSONObject: ["data": cleanArea])
        let decoded = try JSONDecoder().decode(SingleAreaResponse.self, from: json)
        return decoded.data
    }

    private static func describeDecodingError(_ error: Error) -> String {
        if let de = error as? DecodingError {
            switch de {
            case .keyNotFound(let key, let ctx):
                return "keyNotFound(\(key.stringValue)) at path: \(ctx.codingPath.map { $0.stringValue }.joined(separator: "."))"
            case .typeMismatch(let type, let ctx):
                return "typeMismatch(\(type)) at path: \(ctx.codingPath.map { $0.stringValue }.joined(separator: "."))"
            case .valueNotFound(let type, let ctx):
                return "valueNotFound(\(type)) at path: \(ctx.codingPath.map { $0.stringValue }.joined(separator: "."))"
            case .dataCorrupted(let ctx):
                return "dataCorrupted at path: \(ctx.codingPath.map { $0.stringValue }.joined(separator: "."))"
            @unknown default:
                return de.localizedDescription
            }
        }
        return error.localizedDescription
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
    var name: String
    var description: String
    var triggerService: String
    var actionService: String
    var isActive: Bool
    let isPublic: Bool
    let createdAt: String
    var updatedAt: String
    let userID: Int
    var triggerIconURL: String?
    var actionIconURL: String?
    var status: String?
    var triggerType: String?
    var actionType: String?
    var triggerConfig: [String: AnyCodable]?
    var actionConfig: [String: AnyCodable]?
    let conditions: [AnyCodable]?
    var scheduleCron: String?
    var rateLimitPerMin: Int?
    var dedupWindowSec: Int?
    var retryMax: Int?
    var retryBackoffMs: Int?
    var lastRunStatus: String?
    var lastRunAt: String?
    var nextRunAt: String?
    var runCount: Int?
    var lastError: String?
    var dedupKeyTemplate: String?
    let user: AreaUser?
    
    enum CodingKeys: String, CodingKey {
        case id, name, description, status, conditions, user
        case isActive = "is_active"
        case isPublic = "is_public"
        case triggerService = "trigger_service"
        case actionService = "action_service"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case userID = "user_id"
        case triggerIconURL = "trigger_icon_url"
        case actionIconURL = "action_icon_url"
        case triggerType = "trigger_type"
        case actionType = "action_type"
        case triggerConfig = "trigger_config"
        case actionConfig = "action_config"
        case scheduleCron = "schedule_cron"
        case rateLimitPerMin = "rate_limit_per_min"
        case dedupWindowSec = "dedup_window_sec"
        case retryMax = "retry_max"
        case retryBackoffMs = "retry_backoff_ms"
        case lastRunStatus = "last_run_status"
        case lastRunAt = "last_run_at"
        case nextRunAt = "next_run_at"
        case runCount = "run_count"
        case lastError = "last_error"
        case dedupKeyTemplate = "dedup_key_template"
    }
}

// MARK: - Create/Update Models
struct CreateOrUpdateAreaRequest: Codable {
    var name: String?
    var description: String?
    var triggerService: String?
    var triggerType: String?
    var actionService: String?
    var actionType: String?
    var triggerConfig: [String: AnyCodable]?
    var actionConfig: [String: AnyCodable]?
    var isActive: Bool?
}

struct AreaUser: Codable {
    let id: Int
    let email: String
    let firstName: String?
    let lastName: String?
    let createdAt: String
    let updatedAt: String
    let phone: String?
    let birthday: String?
    let gender: String?
    let country: String?
    let lang: String?
    let loginProvider: String?
    let profileImage: String?
    let githubID: String?
    let githubUsername: String?
    let googleID: String?
    let googleEmail: String?
    let facebookID: String?
    let facebookEmail: String?
    let role: String?
    let isActive: Bool?
    
    enum CodingKeys: String, CodingKey {
        case id, email, phone, birthday, gender, country, lang, role
        case firstName = "first_name"
        case lastName = "last_name"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case loginProvider = "login_provider"
        case profileImage = "profile_image"
        case githubID = "github_id"
        case githubUsername = "github_username"
        case googleID = "google_id"
        case googleEmail = "google_email"
        case facebookID = "facebook_id"
        case facebookEmail = "facebook_email"
        case isActive = "is_active"
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
    
    static func from(area: Area) -> AreaTemplate {
        AreaTemplate(
            id: area.id,
            title: area.name,
            subtitle: "Update this AREA",
            description: area.description,
            icon: "gear",
            gradientClass: "gray",
            triggerService: area.triggerService,
            actionService: area.actionService,
            triggerIconURL: area.triggerIconURL,
            actionIconURL: area.actionIconURL,
            isActive: area.isActive
        )
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
