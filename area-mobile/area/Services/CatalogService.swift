//
//  CatalogService.swift
//  area
//
//  Fetches /about.json and exposes services/actions/reactions to the app
//

import Foundation

@MainActor
class CatalogService: ObservableObject {
    static let shared = CatalogService()

    @Published var services: [AboutService] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private init() {}

    func fetchCatalog() async {
        isLoading = true
        errorMessage = nil

        guard let url = URL(string: AppConfig.getAPIEndpoint("/about.json")) else {
            self.errorMessage = "Invalid URL"
            self.isLoading = false
            print("❌ CatalogService: Invalid URL")
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        do {
            print("📡 CatalogService: Fetching catalog from \(url.absoluteString)")
            let (data, response) = try await URLSession.shared.data(for: request)
            
            if let http = response as? HTTPURLResponse {
                print("📡 CatalogService: Response status: \(http.statusCode)")
                if http.statusCode != 200 {
                    let errorBody = String(data: data, encoding: .utf8) ?? "No response body"
                    self.errorMessage = "Failed to load catalog (status: \(http.statusCode))"
                    self.isLoading = false
                    print("❌ CatalogService: HTTP error \(http.statusCode): \(errorBody)")
                    return
                }
            }

            // Print raw JSON for debugging
            if let jsonString = String(data: data, encoding: .utf8) {
                print("📦 CatalogService: Received JSON (first 500 chars): \(String(jsonString.prefix(500)))")
            }

            let decoded = try JSONDecoder().decode(AboutResponse.self, from: data)
            self.services = decoded.server.services
            self.isLoading = false
            
            print("✅ CatalogService: Successfully loaded \(self.services.count) services")
            for service in self.services {
                print("   - \(service.name): \(service.actions.count) actions, \(service.reactions.count) reactions")
            }
        } catch {
            let errorDetails = error.localizedDescription
            self.errorMessage = "Network error: \(errorDetails)"
            self.isLoading = false
            print("❌ CatalogService: Decoding error: \(error)")
            
            // Try to print more details about the error
            if let decodingError = error as? DecodingError {
                switch decodingError {
                case .keyNotFound(let key, let context):
                    print("   Missing key: \(key.stringValue) at \(context.codingPath)")
                case .typeMismatch(let type, let context):
                    print("   Type mismatch: expected \(type) at \(context.codingPath)")
                case .valueNotFound(let type, let context):
                    print("   Value not found: \(type) at \(context.codingPath)")
                case .dataCorrupted(let context):
                    print("   Data corrupted: \(context.debugDescription)")
                @unknown default:
                    print("   Unknown decoding error")
                }
            }
        }
    }
}


