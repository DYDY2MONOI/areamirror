//
//  OAuthConfig.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import Foundation

struct OAuthConfig {
    private static func infoValue(for key: String) -> String? {
        Bundle.main.object(forInfoDictionaryKey: key) as? String
    }

    static var callbackScheme: String {
        infoValue(for: "OAUTH_CALLBACK_SCHEME") ?? "area"
    }

    static var callbackHost: String {
        infoValue(for: "OAUTH_CALLBACK_HOST") ?? "oauth2"
    }

    static var callbackPath: String {
        infoValue(for: "OAUTH_CALLBACK_PATH") ?? "/callback"
    }

    static var callbackURLString: String {
        "\(callbackScheme)://\(callbackHost)\(callbackPath)"
    }

    static func clientID(for providerID: String) -> String {
        let key = "\(providerID.uppercased())_CLIENT_ID"
        switch providerID {
        case "github":
            return infoValue(for: key) ?? "Ov23liQ7GPEEWs0hVzyM"
        case "google":
            return infoValue(for: key) ?? "635456942993-cl4hm8rvh0pcslhjfko23cbgh18o4cdq.apps.googleusercontent.com"
        case "spotify":
            return infoValue(for: key) ?? ""
        case "twitter":
            return infoValue(for: key) ?? "bk1hWG93bjJwOGRqbXh6ZjJpRkQ6MTpjaQ"
        case "facebook":
            return infoValue(for: key) ?? "827003503004798"
        default:
            return infoValue(for: key) ?? ""
        }
    }

    static func providerRedirectURI(for provider: OAuthProvider) -> String {
        let overrideKey = "MOBILE_\(provider.id.uppercased())_REDIRECT_URI"
        if let override = infoValue(for: overrideKey), !override.isEmpty {
            return override
        }

        return "\(AppConfig.baseURL)\(provider.callbackPath)"
    }
}

