//
//  OAuthLoginManager.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import AuthenticationServices
import CryptoKit
import Foundation
import UIKit

enum OAuthLoginError: LocalizedError, Equatable {
    case providerDisabled
    case invalidConfiguration
    case authenticationFailed
    case userCancelled
    case missingCallbackData

    var errorDescription: String? {
        switch self {
        case .providerDisabled:
            return "This authentication provider is currently disabled."
        case .invalidConfiguration:
            return "OAuth configuration is invalid. Please check client identifiers and redirect URIs."
        case .authenticationFailed:
            return "Unable to complete authentication with the provider."
        case .userCancelled:
            return "Authentication was cancelled."
        case .missingCallbackData:
            return "Missing authentication data in callback response."
        }
    }
}

struct OAuthCallbackTokens {
    let providerID: String
    let accessToken: String
    let refreshToken: String?
    let expiresIn: Int
    let tokenType: String
}

struct OAuthLinkCallback {
    let providerID: String
    let authorizationCode: String
    let codeVerifier: String?
}

@available(iOS 13.0, *)
final class OAuthLoginManager: NSObject, ObservableObject {
    private enum Flow {
        case login
        case link
    }

    static let shared = OAuthLoginManager()

    private var session: ASWebAuthenticationSession?
    private var currentFlow: Flow?
    private var currentProvider: OAuthProvider?
    private var currentCodeVerifier: String?
    private var loginCompletion: ((Result<OAuthCallbackTokens, Error>) -> Void)?
    private var linkCompletion: ((Result<OAuthLinkCallback, Error>) -> Void)?

    func startLogin(with provider: OAuthProvider, completion: @escaping (Result<OAuthCallbackTokens, Error>) -> Void) {
        guard provider.isEnabled else {
            completion(.failure(OAuthLoginError.providerDisabled))
            return
        }

        guard !provider.clientId.isEmpty else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            return
        }

        guard loginCompletion == nil && linkCompletion == nil else {
            completion(.failure(OAuthLoginError.authenticationFailed))
            return
        }

        currentFlow = .login
        currentProvider = provider
        loginCompletion = completion

        guard let authorizationURL = buildAuthorizationURL(for: provider, flow: .login) else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            cleanup()
            return
        }

        startSession(with: authorizationURL) { [weak self] result in
            guard let self else { return }
            switch result {
            case .success(let url):
                self.handleCallbackURL(url)
            case .failure(let error):
                self.failCurrentFlow(with: error)
            }
        }
    }

    func startLink(with provider: OAuthProvider, completion: @escaping (Result<OAuthLinkCallback, Error>) -> Void) {
        guard provider.isEnabled else {
            completion(.failure(OAuthLoginError.providerDisabled))
            return
        }

        guard !provider.clientId.isEmpty else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            return
        }

        guard loginCompletion == nil && linkCompletion == nil else {
            completion(.failure(OAuthLoginError.authenticationFailed))
            return
        }

        currentFlow = .link
        currentProvider = provider
        linkCompletion = completion

        guard let authorizationURL = buildAuthorizationURL(for: provider, flow: .link) else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            cleanup()
            return
        }

        startSession(with: authorizationURL) { [weak self] result in
            guard let self else { return }
            switch result {
            case .success(let url):
                self.handleCallbackURL(url)
            case .failure(let error):
                self.failCurrentFlow(with: error)
            }
        }
    }

    @MainActor
    func performLink(with provider: OAuthProvider) async throws -> OAuthLinkCallback {
        try await withCheckedThrowingContinuation { continuation in
            startLink(with: provider) { result in
                continuation.resume(with: result)
            }
        }
    }

    private func startSession(with url: URL, completion: @escaping (Result<URL, Error>) -> Void) {
        let session = ASWebAuthenticationSession(url: url, callbackURLScheme: OAuthConfig.callbackScheme) { [weak self] callbackURL, error in
            guard let self else { return }
            defer { self.cleanup() }

            if let error = error as? ASWebAuthenticationSessionError, error.code == .canceledLogin {
                completion(.failure(OAuthLoginError.userCancelled))
                return
            }

            if let error {
                completion(.failure(error))
                return
            }

            guard let callbackURL else {
                completion(.failure(OAuthLoginError.missingCallbackData))
                return
            }

            completion(.success(callbackURL))
        }

        session.presentationContextProvider = self
        session.prefersEphemeralWebBrowserSession = true
        self.session = session

        if !session.start() {
            completion(.failure(OAuthLoginError.authenticationFailed))
            cleanup()
        }
    }

    private func buildAuthorizationURL(for provider: OAuthProvider, flow: Flow) -> URL? {
        guard var components = URLComponents(string: provider.authorizationEndpoint) else {
            return nil
        }

        var queryItems = [
            URLQueryItem(name: "client_id", value: provider.clientId),
            URLQueryItem(name: "redirect_uri", value: provider.redirectURI),
            URLQueryItem(name: "response_type", value: provider.responseType),
            URLQueryItem(name: "scope", value: provider.scopeValue)
        ]

        currentCodeVerifier = nil
        var stateComponents: [String] = ["mobile"]

        switch flow {
        case .login:
            if provider.requiresPKCE {
                let verifier = PKCEHelper.generateCodeVerifier()
                currentCodeVerifier = verifier
                let challenge = PKCEHelper.codeChallenge(for: verifier)
                queryItems.append(URLQueryItem(name: "code_challenge", value: challenge))
                queryItems.append(URLQueryItem(name: "code_challenge_method", value: "S256"))
                stateComponents.append(verifier)
            }
        case .link:
            stateComponents.append("link")
            if provider.requiresPKCE {
                let verifier = PKCEHelper.generateCodeVerifier()
                currentCodeVerifier = verifier
                let challenge = PKCEHelper.codeChallenge(for: verifier)
                queryItems.append(URLQueryItem(name: "code_challenge", value: challenge))
                queryItems.append(URLQueryItem(name: "code_challenge_method", value: "S256"))
                stateComponents.append(verifier)
            }
        }

        let state = stateComponents.joined(separator: ":")
        queryItems.append(URLQueryItem(name: "state", value: state))

        provider.additionalParameters.forEach { key, value in
            queryItems.append(URLQueryItem(name: key, value: value))
        }

        components.queryItems = queryItems
        return components.url
    }

    private func handleCallbackURL(_ url: URL) {
        guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let items = components.queryItems else {
            failCurrentFlow(with: OAuthLoginError.missingCallbackData)
            return
        }

        var data: [String: String] = [:]
        items.forEach { data[$0.name] = $0.value }

        let mode = data["mode"]

        if mode == "link" || currentFlow == .link {
            guard let providerID = data["provider"] ?? currentProvider?.id,
                  let code = data["code"] else {
                failCurrentFlow(with: OAuthLoginError.missingCallbackData)
                return
            }

            let verifier = data["code_verifier"] ?? currentCodeVerifier
            let callback = OAuthLinkCallback(providerID: providerID, authorizationCode: code, codeVerifier: verifier)
            completeLink(with: callback)
            return
        }

        guard
            let providerID = data["provider"] ?? currentProvider?.id,
            let accessToken = data["access_token"],
            let expiresString = data["expires_in"],
            let expiresIn = Int(expiresString),
            let tokenType = data["token_type"]
        else {
            failCurrentFlow(with: OAuthLoginError.missingCallbackData)
            return
        }

        let refreshToken = data["refresh_token"]
        let tokens = OAuthCallbackTokens(
            providerID: providerID,
            accessToken: accessToken,
            refreshToken: refreshToken,
            expiresIn: expiresIn,
            tokenType: tokenType
        )

        completeLogin(with: tokens)
    }

    private func failCurrentFlow(with error: Error) {
        switch currentFlow {
        case .login:
            if let completion = loginCompletion {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
            }
        case .link:
            if let completion = linkCompletion {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
            }
        case .none:
            break
        }
    }

    private func completeLogin(with tokens: OAuthCallbackTokens) {
        guard let completion = loginCompletion else { return }
        DispatchQueue.main.async {
            completion(.success(tokens))
        }
    }

    private func completeLink(with callback: OAuthLinkCallback) {
        guard let completion = linkCompletion else { return }
        DispatchQueue.main.async {
            completion(.success(callback))
        }
    }

    private func cleanup() {
        session = nil
        currentFlow = nil
        currentProvider = nil
        currentCodeVerifier = nil
        loginCompletion = nil
        linkCompletion = nil
    }
}

@available(iOS 13.0, *)
extension OAuthLoginManager: ASWebAuthenticationPresentationContextProviding {
    func presentationAnchor(for session: ASWebAuthenticationSession) -> ASPresentationAnchor {
        if let scene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
           let window = scene.windows.first(where: { $0.isKeyWindow }) {
            return window
        }

        return UIWindow()
    }
}

private enum PKCEHelper {
    static func generateCodeVerifier() -> String {
        let characters = Array("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")
        var verifier = ""
        verifier.reserveCapacity(64)

        for _ in 0..<64 {
            verifier.append(characters.randomElement() ?? "a")
        }

        return verifier
    }

    static func codeChallenge(for verifier: String) -> String {
        guard let data = verifier.data(using: .utf8) else {
            return verifier
        }

        let hash = SHA256.hash(data: data)
        return base64URLEncode(Data(hash))
    }

    private static func base64URLEncode(_ data: Data) -> String {
        data.base64EncodedString()
            .replacingOccurrences(of: "+", with: "-")
            .replacingOccurrences(of: "/", with: "_")
            .replacingOccurrences(of: "=", with: "")
    }
}

