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
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)
            if let http = response as? HTTPURLResponse, http.statusCode != 200 {
                self.errorMessage = "Failed to load catalog (status: \(http.statusCode))"
                self.isLoading = false
                return
            }

            let decoded = try JSONDecoder().decode(AboutResponse.self, from: data)
            self.services = decoded.server.services
            self.isLoading = false
        } catch {
            self.errorMessage = "Network error: \(error.localizedDescription)"
            self.isLoading = false
        }
    }
}


