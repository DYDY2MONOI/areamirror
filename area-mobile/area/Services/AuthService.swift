//
//  AuthService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import Combine
import UIKit

class AuthService: ObservableObject {
    static let shared = AuthService()
    
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let baseURL = AppConfig.baseURL
    private var cancellables = Set<AnyCancellable>()
    
    private init() {
        checkAuthStatus()
    }
    
    
    private func checkAuthStatus() {
        if let token = UserDefaults.standard.string(forKey: "authToken"), !token.isEmpty {
            fetchProfile()
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
    
    func removeToken() {
        UserDefaults.standard.removeObject(forKey: "authToken")
    }
    
    func fetchProfile() {
        guard let url = URL(string: "\(baseURL)/profile") else {
            errorMessage = "Invalid URL"
            isAuthenticated = false
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                if let error = error {
                    self?.errorMessage = "Network error: \(error.localizedDescription)"
                    self?.isAuthenticated = false
                    return
                }
                
                guard let data = data else {
                    self?.errorMessage = "No data received"
                    self?.isAuthenticated = false
                    return
                }
                
                if let httpResponse = response as? HTTPURLResponse {
                    if httpResponse.statusCode == 200 {
                        do {
                            let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                            self?.currentUser = profileResponse.user
                            self?.isAuthenticated = true
                            self?.errorMessage = nil
                        } catch {
                            self?.errorMessage = "Error decoding response"
                            self?.isAuthenticated = false
                        }
                    } else {
                        do {
                            let errorResponse = try JSONDecoder().decode([String: String].self, from: data)
                            self?.errorMessage = errorResponse["error"] ?? "Profile fetch error"
                        } catch {
                            self?.errorMessage = "Profile fetch error"
                        }
                        self?.isAuthenticated = false
                        self?.removeToken()
                    }
                }
            }
        }.resume()
    }
    
    func updateProfile(firstName: String?, lastName: String?, phone: String?, country: String?, currentPassword: String?, newPassword: String?) {
        isLoading = true
        errorMessage = nil
        
        let updateRequest = ProfileUpdateRequest(
            firstName: firstName,
            lastName: lastName,
            phone: phone,
            country: country,
            currentPassword: currentPassword,
            newPassword: newPassword
        )
        
        guard let url = URL(string: "\(baseURL)/profile") else {
            errorMessage = "Invalid URL"
            isLoading = false
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        do {
            request.httpBody = try JSONEncoder().encode(updateRequest)
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
                            let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                            self?.currentUser = profileResponse.user
                            self?.errorMessage = nil
                        } catch {
                            self?.errorMessage = "Error decoding response: \(error.localizedDescription)"
                        }
                    } else {
                        do {
                            let errorResponse = try JSONDecoder().decode([String: String].self, from: data)
                            self?.errorMessage = errorResponse["error"] ?? "Profile update error"
                        } catch {
                            let responseString = String(data: data, encoding: .utf8) ?? "Unknown error"
                            self?.errorMessage = "Profile update error: \(responseString)"
                        }
                    }
                }
            }
        }.resume()
    }
    
    func uploadProfileImage(_ image: UIImage) {
        isLoading = true
        errorMessage = nil
        
        guard let url = URL(string: "\(baseURL)/profile/image") else {
            errorMessage = "Invalid URL"
            isLoading = false
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        
        if let token = getAuthToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        
        let boundary = UUID().uuidString
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
        
        var body = Data()
        
        // Add image data
        if let imageData = image.jpegData(compressionQuality: 0.8) {
            body.append("--\(boundary)\r\n".data(using: .utf8)!)
            body.append("Content-Disposition: form-data; name=\"image\"; filename=\"profile.jpg\"\r\n".data(using: .utf8)!)
            body.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
            body.append(imageData)
            body.append("\r\n".data(using: .utf8)!)
        }
        
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)
        request.httpBody = body
        
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
                            let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                            self?.currentUser = profileResponse.user
                            self?.errorMessage = nil
                        } catch {
                            self?.errorMessage = "Error decoding response: \(error.localizedDescription)"
                        }
                    } else {
                        do {
                            let errorResponse = try JSONDecoder().decode([String: String].self, from: data)
                            self?.errorMessage = errorResponse["error"] ?? "Image upload error"
                        } catch {
                            let responseString = String(data: data, encoding: .utf8) ?? "Unknown error"
                            self?.errorMessage = "Image upload error: \(responseString)"
                        }
                    }
                }
            }
        }.resume()
    }
}
