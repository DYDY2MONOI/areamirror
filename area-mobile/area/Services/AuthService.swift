//
//  AuthService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import Combine
import UIKit

enum AuthServiceError: LocalizedError {
    case missingRefreshToken
    case invalidResponse
    case tokenRefreshFailed
    case unsupportedPlatform
    case invalidProvider
    case custom(String)

    var errorDescription: String? {
        switch self {
        case .missingRefreshToken:
            return "Missing refresh token. Please log in again."
        case .invalidResponse:
            return "Received invalid response from server."
        case .tokenRefreshFailed:
            return "Unable to refresh authentication token."
        case .unsupportedPlatform:
            return "This feature requires a newer version of iOS."
        case .invalidProvider:
            return "Unsupported authentication provider."
        case .custom(let message):
            return message
        }
    }
}

class AuthService: ObservableObject {
    static let shared = AuthService()
    
    @Published var isAuthenticated = false
    @Published var currentUser: User?
    @Published var isLoading = false
    @Published var errorMessage: String?
    
    private let baseURL = AppConfig.baseURL
    private let accessTokenKey = "oauth_access_token"
    private let refreshTokenKey = "oauth_refresh_token"
    private let tokenExpiryKey = "oauth_token_expiry"
    private let tokenTypeKey = "oauth_token_type"

    private var accessToken: String?
    private var refreshToken: String?
    private var tokenExpiry: Date?
    private var tokenType: String = "Bearer"
    
    private init() {
        loadStoredTokens()
    }
    
    
    private func loadStoredTokens() {
        let defaults = UserDefaults.standard

        if let storedAccessToken = defaults.string(forKey: accessTokenKey) {
            accessToken = storedAccessToken
        }

        if let storedRefreshToken = defaults.string(forKey: refreshTokenKey) {
            refreshToken = storedRefreshToken
        }

        tokenType = defaults.string(forKey: tokenTypeKey) ?? "Bearer"

        let expiryTimestamp = defaults.double(forKey: tokenExpiryKey)
        if expiryTimestamp > 0 {
            tokenExpiry = Date(timeIntervalSince1970: expiryTimestamp)
        }

        guard accessToken != nil else {
            isAuthenticated = false
            return
        }

        if shouldRefreshToken {
            refreshAccessToken { _ in }
        } else {
            fetchProfile()
        }
    }

    private var shouldRefreshToken: Bool {
        guard let expiry = tokenExpiry else { return false }
        return Date() >= expiry.addingTimeInterval(-60)
    }
    
    
    func login(email: String, password: String) {
        isLoading = true
        errorMessage = nil

        let loginRequest = LoginRequest(email: email, password: password)

        guard let url = URL(string: "\(baseURL)/mobile/oauth2/login") else {
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
                guard let self else { return }

                self.isLoading = false

                if let error = error {
                    self.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }

                guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                    self.errorMessage = "No data received"
                    return
                }

                if httpResponse.statusCode == 200 {
                    do {
                        let authResponse = try JSONDecoder().decode(OAuthTokenResponse.self, from: data)
                        self.handleOAuthTokenResponse(authResponse)
                    } catch {
                        self.errorMessage = "Error decoding response"
                    }
                } else {
                    self.errorMessage = AuthService.decodeErrorMessage(from: data) ?? "Login error"
                }
            }
        }.resume()
    }

    func login(with provider: OAuthProvider) {
        guard #available(iOS 13.0, *) else {
            errorMessage = "OAuth login requires iOS 13 or later."
            return
        }

        isLoading = true
        errorMessage = nil

        OAuthLoginManager.shared.startLogin(with: provider) { [weak self] result in
            DispatchQueue.main.async {
                guard let self else { return }

                self.isLoading = false

                switch result {
                case .success(let tokens):
                    self.handleOAuthCallback(tokens: tokens)
                case .failure(let error):
                    if let oauthError = error as? OAuthLoginError, oauthError == .userCancelled {
                        self.errorMessage = nil
                    } else {
                        self.errorMessage = error.localizedDescription
                    }
                }
            }
        }
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
                guard let self else { return }

                self.isLoading = false

                if let error = error {
                    self.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }

                guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                    self.errorMessage = "No data received"
                    return
                }

                if httpResponse.statusCode == 201 {
                    self.login(email: email, password: password)
                } else {
                    self.errorMessage = AuthService.decodeErrorMessage(from: data) ?? "Registration error"
                }
            }
        }.resume()
    }
    
    func logout() {
        clearStoredTokens()
        isAuthenticated = false
        currentUser = nil
        errorMessage = nil
    }
    
    func fetchProfile() {
        ensureValidToken { [weak self] result in
            guard let self else { return }

            switch result {
            case .success(let token):
                self.performFetchProfile(with: token, retryingAfterRefresh: true)
            case .failure(let error):
                DispatchQueue.main.async {
                    self.errorMessage = error.localizedDescription
                    self.isAuthenticated = false
                }
            }
        }
    }

    private func performFetchProfile(with token: String, retryingAfterRefresh: Bool) {
        guard let url = URL(string: "\(baseURL)/mobile/oauth2/me") else {
            DispatchQueue.main.async {
                self.errorMessage = "Invalid URL"
                self.isAuthenticated = false
            }
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("\(tokenType) \(token)", forHTTPHeaderField: "Authorization")

        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                guard let self else { return }

                if let error = error {
                    self.errorMessage = "Network error: \(error.localizedDescription)"
                    self.isAuthenticated = false
                    return
                }

                guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                    self.errorMessage = "No data received"
                    self.isAuthenticated = false
                    return
                }

                switch httpResponse.statusCode {
                case 200:
                    do {
                        let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                        self.currentUser = profileResponse.user
                        self.isAuthenticated = true
                        self.errorMessage = nil
                    } catch {
                        self.errorMessage = "Error decoding profile"
                        self.isAuthenticated = false
                    }

                case 401 where retryingAfterRefresh:
                    self.refreshAccessToken { result in
                        switch result {
                        case .success(let newToken):
                            self.performFetchProfile(with: newToken, retryingAfterRefresh: false)
                        case .failure(let refreshError):
                            DispatchQueue.main.async {
                                self.errorMessage = refreshError.localizedDescription
                                self.isAuthenticated = false
                            }
                        }
                    }

                default:
                    self.errorMessage = AuthService.decodeErrorMessage(from: data) ?? "Profile fetch error"
                    self.isAuthenticated = false
                }
            }
        }.resume()
    }

    private func handleOAuthTokenResponse(_ response: OAuthTokenResponse) {
        storeTokens(
            accessToken: response.accessToken,
            refreshToken: response.refreshToken,
            expiresIn: response.expiresIn,
            tokenType: response.tokenType
        )

        currentUser = response.user
        isAuthenticated = true
        errorMessage = nil
    }

    private func handleOAuthCallback(tokens: OAuthCallbackTokens) {
        storeTokens(
            accessToken: tokens.accessToken,
            refreshToken: tokens.refreshToken,
            expiresIn: tokens.expiresIn,
            tokenType: tokens.tokenType
        )

        fetchProfile()
    }

    private func storeTokens(accessToken: String, refreshToken: String?, expiresIn: Int, tokenType: String) {
        self.accessToken = accessToken
        if let refreshToken, !refreshToken.isEmpty {
            self.refreshToken = refreshToken
        }
        self.tokenType = tokenType

        let expiryDate = Date().addingTimeInterval(TimeInterval(expiresIn))
        tokenExpiry = expiryDate

        let defaults = UserDefaults.standard
        defaults.set(accessToken, forKey: accessTokenKey)
        defaults.set(self.refreshToken, forKey: refreshTokenKey)
        defaults.set(expiryDate.timeIntervalSince1970, forKey: tokenExpiryKey)
        defaults.set(tokenType, forKey: tokenTypeKey)
    }

    private func clearStoredTokens() {
        accessToken = nil
        refreshToken = nil
        tokenExpiry = nil
        tokenType = "Bearer"

        let defaults = UserDefaults.standard
        defaults.removeObject(forKey: accessTokenKey)
        defaults.removeObject(forKey: refreshTokenKey)
        defaults.removeObject(forKey: tokenExpiryKey)
        defaults.removeObject(forKey: tokenTypeKey)
    }

    func getAuthToken() -> String? {
        accessToken
    }

    func authorizationHeader() -> String? {
        guard let token = accessToken else { return nil }
        return "\(tokenType) \(token)"
    }

    @MainActor
    func getValidToken() async throws -> String {
        try await withCheckedThrowingContinuation { continuation in
            ensureValidToken { result in
                continuation.resume(with: result)
            }
        }
    }

    @MainActor
    func getValidAuthorizationHeader() async throws -> String {
        let token = try await getValidToken()
        return "\(tokenType) \(token)"
    }

    private func ensureValidToken(completion: @escaping (Result<String, Error>) -> Void) {
        if let token = accessToken, !shouldRefreshToken {
            completion(.success(token))
            return
        }

        refreshAccessToken(completion: completion)
    }

    private func refreshAccessToken(completion: @escaping (Result<String, Error>) -> Void) {
        guard let refreshToken = refreshToken else {
            clearStoredTokens()
            completion(.failure(AuthServiceError.missingRefreshToken))
            return
        }

        guard let url = URL(string: "\(baseURL)/mobile/oauth2/refresh") else {
            completion(.failure(AuthServiceError.invalidResponse))
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let payload = RefreshTokenRequest(refreshToken: refreshToken)

        do {
            request.httpBody = try JSONEncoder().encode(payload)
        } catch {
            completion(.failure(error))
            return
        }

        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            guard let self else { return }

            if let error = error {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
                return
            }

            guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                DispatchQueue.main.async {
                    completion(.failure(AuthServiceError.invalidResponse))
                }
                return
            }

            if httpResponse.statusCode == 200 {
                do {
                    let refreshResponse = try JSONDecoder().decode(RefreshTokenResponse.self, from: data)
                    self.storeTokens(
                        accessToken: refreshResponse.accessToken,
                        refreshToken: self.refreshToken,
                        expiresIn: refreshResponse.expiresIn,
                        tokenType: refreshResponse.tokenType
                    )

                    DispatchQueue.main.async {
                        completion(.success(refreshResponse.accessToken))
                    }
                } catch {
                    DispatchQueue.main.async {
                        completion(.failure(error))
                    }
                }
            } else {
                DispatchQueue.main.async {
                    self.clearStoredTokens()
                    completion(.failure(AuthServiceError.tokenRefreshFailed))
                }
            }
        }.resume()
    }

    @MainActor
    func linkAccount(_ provider: OAuthProvider) async throws {
        guard #available(iOS 13.0, *) else {
            throw AuthServiceError.unsupportedPlatform
        }

        let linkResult = try await OAuthLoginManager.shared.performLink(with: provider)

        guard linkResult.providerID == provider.id else {
            throw AuthServiceError.invalidProvider
        }

        try await sendLinkRequest(providerID: provider.id, code: linkResult.authorizationCode, codeVerifier: linkResult.codeVerifier)
        fetchProfile()
    }

    @MainActor
    func unlinkAccount(_ provider: OAuthProvider) async throws {
        try await sendUnlinkRequest(providerID: provider.id)
        fetchProfile()
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
        
        guard let body = try? JSONEncoder().encode(updateRequest) else {
            errorMessage = "Error encoding data"
            isLoading = false
            return
        }

        ensureValidToken { [weak self] result in
            guard let self else { return }

            switch result {
            case .success(let token):
                self.performProfileUpdate(with: token, body: body)
            case .failure(let error):
                DispatchQueue.main.async {
                    self.isLoading = false
                    self.errorMessage = error.localizedDescription
                }
            }
        }
    }

    private func performProfileUpdate(with token: String, body: Data) {
        guard let url = URL(string: "\(baseURL)/profile") else {
            DispatchQueue.main.async {
                self.isLoading = false
                self.errorMessage = "Invalid URL"
            }
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("\(tokenType) \(token)", forHTTPHeaderField: "Authorization")
        request.httpBody = body

        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                guard let self else { return }

                self.isLoading = false

                if let error = error {
                    self.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }

                guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                    self.errorMessage = "No data received"
                    return
                }

                if httpResponse.statusCode == 200 {
                    do {
                        let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                        self.currentUser = profileResponse.user
                        self.errorMessage = nil
                    } catch {
                        self.errorMessage = "Error decoding response: \(error.localizedDescription)"
                    }
                } else {
                    let fallback = String(data: data, encoding: .utf8) ?? "Profile update error"
                    self.errorMessage = AuthService.decodeErrorMessage(from: data) ?? fallback
                }
            }
        }.resume()
    }
    
    func uploadProfileImage(_ image: UIImage) {
        isLoading = true
        errorMessage = nil
        
        guard let imageData = image.jpegData(compressionQuality: 0.8) else {
            errorMessage = "Unable to process image"
            isLoading = false
            return
        }

        ensureValidToken { [weak self] result in
            guard let self else { return }

            switch result {
            case .success(let token):
                self.performProfileImageUpload(with: token, imageData: imageData)
            case .failure(let error):
                DispatchQueue.main.async {
                    self.isLoading = false
                    self.errorMessage = error.localizedDescription
                }
            }
        }
    }

    private func performProfileImageUpload(with token: String, imageData: Data) {
        guard let url = URL(string: "\(baseURL)/profile/image") else {
            DispatchQueue.main.async {
                self.isLoading = false
                self.errorMessage = "Invalid URL"
            }
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("\(tokenType) \(token)", forHTTPHeaderField: "Authorization")

        let boundary = UUID().uuidString
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")

        var body = Data()
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"image\"; filename=\"profile.jpg\"\r\n".data(using: .utf8)!)
        body.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
        body.append(imageData)
        body.append("\r\n".data(using: .utf8)!)
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)
        request.httpBody = body

        URLSession.shared.dataTask(with: request) { [weak self] data, response, error in
            DispatchQueue.main.async {
                guard let self else { return }

                self.isLoading = false

                if let error = error {
                    self.errorMessage = "Network error: \(error.localizedDescription)"
                    return
                }

                guard let data = data, let httpResponse = response as? HTTPURLResponse else {
                    self.errorMessage = "No data received"
                    return
                }

                if httpResponse.statusCode == 200 {
                    do {
                        let profileResponse = try JSONDecoder().decode(ProfileResponse.self, from: data)
                        self.currentUser = profileResponse.user
                        self.errorMessage = nil
                    } catch {
                        self.errorMessage = "Error decoding response: \(error.localizedDescription)"
                    }
                } else {
                    let fallback = String(data: data, encoding: .utf8) ?? "Image upload error"
                    self.errorMessage = AuthService.decodeErrorMessage(from: data) ?? fallback
                }
            }
        }.resume()
    }

@MainActor
func isProviderLinked(_ providerID: String) -> Bool {
    guard let user = currentUser else { return false }

    switch providerID {
    case "github":
        return !(user.githubUsername?.isEmpty ?? true)
    case "google":
        return !(user.googleID?.isEmpty ?? true)
    case "facebook":
        return !(user.facebookID?.isEmpty ?? true)
    case "spotify":
        return !(user.spotifyID?.isEmpty ?? true)
    case "twitter":
        return !(user.twitterUsername?.isEmpty ?? true)
    default:
        return false
    }
}

@MainActor
func linkedDetail(for providerID: String) -> String? {
    guard let user = currentUser else { return nil }

    switch providerID {
    case "github":
        return user.githubUsername
    case "google":
        return user.googleEmail
    case "facebook":
        return user.facebookEmail
    case "spotify":
        return user.spotifyEmail
    case "twitter":
        return user.twitterUsername
    default:
        return nil
    }
}

@MainActor
private func sendLinkRequest(providerID: String, code: String, codeVerifier: String?) async throws {
    guard let endpoint = linkEndpoint(for: providerID) else {
        throw AuthServiceError.invalidProvider
    }

    var payload: [String: Any] = ["code": code]
    if providerID == "twitter" {
        guard let verifier = codeVerifier, !verifier.isEmpty else {
            throw AuthServiceError.custom("Missing verification data for Twitter linking.")
        }
        payload["code_verifier"] = verifier
    }

    let body = try JSONSerialization.data(withJSONObject: payload, options: [])
    _ = try await performAuthorizedRequest(path: endpoint, method: "POST", body: body)
}

@MainActor
private func sendUnlinkRequest(providerID: String) async throws {
    guard let endpoint = unlinkEndpoint(for: providerID) else {
        throw AuthServiceError.invalidProvider
    }

    _ = try await performAuthorizedRequest(path: endpoint, method: "DELETE")
}

@MainActor
private func performAuthorizedRequest(path: String, method: String, body: Data? = nil) async throws -> Data {
    let header = try await getValidAuthorizationHeader()

    guard let url = URL(string: AppConfig.getAPIEndpoint(path)) else {
        throw AuthServiceError.invalidResponse
    }

    var request = URLRequest(url: url)
    request.httpMethod = method
    request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.setValue(header, forHTTPHeaderField: "Authorization")
    request.httpBody = body

    let (data, response) = try await URLSession.shared.data(for: request)

    guard let httpResponse = response as? HTTPURLResponse else {
        throw AuthServiceError.invalidResponse
    }

    guard (200...299).contains(httpResponse.statusCode) else {
        if let message = AuthService.decodeErrorMessage(from: data) {
            throw AuthServiceError.custom(message)
        }
        throw AuthServiceError.invalidResponse
    }

    return data
}

private func linkEndpoint(for providerID: String) -> String? {
    switch providerID {
    case "github":
        return "/profile/github/link"
    case "google":
        return "/profile/google/link"
    case "facebook":
        return "/profile/facebook/link"
    case "spotify":
        return "/profile/spotify/link"
    case "twitter":
        return "/profile/twitter/link"
    default:
        return nil
    }
}

private func unlinkEndpoint(for providerID: String) -> String? {
    switch providerID {
    case "github":
        return "/profile/github/unlink"
    case "google":
        return "/profile/google/unlink"
    case "facebook":
        return "/profile/facebook/unlink"
    case "spotify":
        return "/profile/spotify/unlink"
    case "twitter":
        return "/profile/twitter/unlink"
    default:
        return nil
    }
}

private static func decodeErrorMessage(from data: Data) -> String? {
    (try? JSONDecoder().decode([String: String].self, from: data))?["error"]
}
}
