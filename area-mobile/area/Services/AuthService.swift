//
//  AuthService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import Combine

class AuthService: ObservableObject {
    static let shared = AuthService()
    
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let baseURL = "http://localhost:8080"
    private var cancellables = Set<AnyCancellable>()
    
    private init() {
        checkAuthStatus()
    }
    
    
    private func checkAuthStatus() {
        if let token = UserDefaults.standard.string(forKey: "authToken"), !token.isEmpty {
            isAuthenticated = true
        }
    }
    
    
    func login(email: String, password: String) {
        isLoading = true
        errorMessage = nil
        
        let loginRequest = LoginRequest(email: email, password: password)
        
        guard let url = URL(string: "\(baseURL)/login") else {
            errorMessage = "Invalid URL"
            isLoading = false
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(loginRequest)
        } catch {
            errorMessage = "Error encoding data"
            isLoading = false
            return
        }
        
        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                self?.isLoading = false
                
                if let error = error {
                    self?.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }
                
                guard let data = data else {
                    self?.errorMessage = "No data received"
                    return
                }
                
                if let httpResponse = response as? HTTPURLResponse {
                    if httpResponse.statusCode == 200 {
                        do {
                            let authResponse = try JSONDecoder().decode(AuthResponse.self, from: data)
                            self?.handleSuccessfulAuth(authResponse)
                        } catch {
                            self?.errorMessage = "Error decoding response"
                        }
                    } else {
                        do {
                            let errorResponse = try JSONDecoder().decode([String: String].self, from: data)
                            self?.errorMessage = errorResponse["error"] ?? "Login error"
                        } catch {
                            self?.errorMessage = "Invalid email or password"
                        }
                    }
                }
            }
        }.resume()
    }
    
    func register(email: String, password: String, firstName: String?, lastName: String?) {
        isLoading = true
        errorMessage = nil
        
        let registerRequest = RegisterRequest(
            email: email,
            password: password,
            firstName: firstName,
            lastName: lastName
        )
        
        guard let url = URL(string: "\(baseURL)/register") else {
            errorMessage = "Invalid URL"
            isLoading = false
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(registerRequest)
        } catch {
            errorMessage = "Error encoding data"
            isLoading = false
            return
        }
        
        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                self?.isLoading = false
                
                if let error = error {
                    self?.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }
                
                guard let data = data else {
                    self?.errorMessage = "No data received"
                    return
                }
                
                if let httpResponse = response as? HTTPURLResponse {
                    if httpResponse.statusCode == 201 {
                        do {
                            let authResponse = try JSONDecoder().decode(AuthResponse.self, from: data)
                            self?.handleSuccessfulAuth(authResponse)
                        } catch {
                            self?.errorMessage = "Error decoding response"
                        }
                    } else {
                        do {
                            let errorResponse = try JSONDecoder().decode([String: String].self, from: data)
                            self?.errorMessage = errorResponse["error"] ?? "Registration error"
                        } catch {
                            self?.errorMessage = "Registration error"
                        }
                    }
                }
            }
        }.resume()
    }
    
    func logout() {
        UserDefaults.standard.removeObject(forKey: "authToken")
        isAuthenticated = false
        currentUser = nil
        errorMessage = nil
    }
    
    private func handleSuccessfulAuth(_ authResponse: AuthResponse) {
        UserDefaults.standard.set(authResponse.token, forKey: "authToken")
        currentUser = authResponse.user
        isAuthenticated = true
        errorMessage = nil
    }
    
    func getAuthToken() -> String? {
        return UserDefaults.standard.string(forKey: "authToken")
    }
}
