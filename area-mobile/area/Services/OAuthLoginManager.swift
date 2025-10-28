//
//  OAuthLoginManager.swift
//  area
//
//  Created by GPT-5 Codex on 28/10/2025.
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

@available(iOS 13.0, *)
final class OAuthLoginManager: NSObject, ObservableObject {
    static let shared = OAuthLoginManager()

    private var session: ASWebAuthenticationSession?
    private var completionHandler: ((Result<OAuthCallbackTokens, Error>) -> Void)?

    func startLogin(with provider: OAuthProvider, completion: @escaping (Result<OAuthCallbackTokens, Error>) -> Void) {
        guard provider.isEnabled else {
            completion(.failure(OAuthLoginError.providerDisabled))
            return
        }

        guard !provider.clientId.isEmpty else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            return
        }

        guard let authorizationURL = buildAuthorizationURL(for: provider) else {
            completion(.failure(OAuthLoginError.invalidConfiguration))
            return
        }

        completionHandler = completion

        let session = ASWebAuthenticationSession(url: authorizationURL, callbackURLScheme: OAuthConfig.callbackScheme) { [weak self] callbackURL, error in
            guard let self else { return }

            defer { self.cleanup() }

            if let error = error as? ASWebAuthenticationSessionError, error.code == .canceledLogin {
                completion(.failure(OAuthLoginError.userCancelled))
                return
            } else if let error {
                completion(.failure(error))
                return
            }

            guard let callbackURL else {
                completion(.failure(OAuthLoginError.missingCallbackData))
                return
            }

            completion(self.handleCallbackURL(callbackURL))
        }

        session.presentationContextProvider = self
        session.prefersEphemeralWebBrowserSession = true
        self.session = session

        if !session.start() {
            completion(.failure(OAuthLoginError.authenticationFailed))
            cleanup()
        }
    }

    private func buildAuthorizationURL(for provider: OAuthProvider) -> URL? {
        guard var components = URLComponents(string: provider.authorizationEndpoint) else {
            return nil
        }

        var queryItems = [
            URLQueryItem(name: "client_id", value: provider.clientId),
            URLQueryItem(name: "redirect_uri", value: provider.redirectURI),
            URLQueryItem(name: "response_type", value: provider.responseType),
            URLQueryItem(name: "scope", value: provider.scopeValue)
        ]

        let state: String
        if provider.requiresPKCE {
            let verifier = PKCEHelper.generateCodeVerifier()
            let challenge = PKCEHelper.codeChallenge(for: verifier)
            queryItems.append(URLQueryItem(name: "code_challenge", value: challenge))
            queryItems.append(URLQueryItem(name: "code_challenge_method", value: "S256"))
            state = "mobile:\(verifier)"
        } else {
            state = "mobile"
        }

        queryItems.append(URLQueryItem(name: "state", value: state))

        provider.additionalParameters.forEach { key, value in
            queryItems.append(URLQueryItem(name: key, value: value))
        }

        components.queryItems = queryItems
        return components.url
    }

    private func handleCallbackURL(_ url: URL) -> Result<OAuthCallbackTokens, Error> {
        guard let components = URLComponents(url: url, resolvingAgainstBaseURL: false),
              let queryItems = components.queryItems else {
            return .failure(OAuthLoginError.missingCallbackData)
        }

        var data: [String: String] = [:]
        queryItems.forEach { data[$0.name] = $0.value }

        guard
            let providerID = data["provider"],
            let accessToken = data["access_token"],
            let expiresString = data["expires_in"],
            let expiresIn = Int(expiresString),
            let tokenType = data["token_type"]
        else {
            return .failure(OAuthLoginError.missingCallbackData)
        }

        let refreshToken = data["refresh_token"]

        let tokens = OAuthCallbackTokens(
            providerID: providerID,
            accessToken: accessToken,
            refreshToken: refreshToken,
            expiresIn: expiresIn,
            tokenType: tokenType
        )

        return .success(tokens)
    }

    private func cleanup() {
        session = nil
        completionHandler = nil
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

