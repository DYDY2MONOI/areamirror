//
//  OAuthProvider.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import SwiftUI

struct OAuthProvider: Identifiable, Hashable {
    enum ScopeSeparator: String {
        case space = " "
        case comma = ","
    }

    let id: String
    let name: String
    let iconName: String
    let color: Color
    let description: String
    let authorizationEndpoint: String
    let callbackPath: String
    let scopes: [String]
    let scopeSeparator: ScopeSeparator
    let responseType: String
    let isEnabled: Bool
    let additionalParameters: [String: String]
    let requiresPKCE: Bool

    var clientId: String {
        OAuthConfig.clientID(for: id)
    }

    var redirectURI: String {
        OAuthConfig.providerRedirectURI(for: self)
    }

    var scopeValue: String {
        scopes.joined(separator: scopeSeparator.rawValue)
    }

    static let availableProviders: [OAuthProvider] = [
        OAuthProvider(
            id: "github",
            name: "GitHub",
            iconName: "github",
            color: Color(hex: "#24292e"),
            description: "Connect your GitHub account to access repositories, issues, and pull requests",
            authorizationEndpoint: "https://github.com/login/oauth/authorize",
            callbackPath: "/oauth2/github/callback",
            scopes: ["user:email"],
            scopeSeparator: .comma,
            responseType: "code",
            isEnabled: true,
            additionalParameters: [
                "allow_signup": "true"
            ],
            requiresPKCE: false
        ),
        OAuthProvider(
            id: "google",
            name: "Google",
            iconName: "gmail",
            color: Color(hex: "#4285f4"),
            description: "Connect your Google account to access Gmail, Calendar, and Drive",
            authorizationEndpoint: "https://accounts.google.com/o/oauth2/v2/auth",
            callbackPath: "/oauth2/google/callback",
            scopes: [
                "openid",
                "email",
                "profile",
                "https://www.googleapis.com/auth/gmail.send",
                "https://www.googleapis.com/auth/gmail.readonly",
                "https://www.googleapis.com/auth/drive.readonly",
                "https://www.googleapis.com/auth/drive.metadata.readonly"
            ],
            scopeSeparator: .space,
            responseType: "code",
            isEnabled: true,
            additionalParameters: [
                "access_type": "offline",
                "prompt": "consent",
                "include_granted_scopes": "true"
            ],
            requiresPKCE: false
        ),
        OAuthProvider(
            id: "facebook",
            name: "Facebook",
            iconName: "facebook",
            color: Color(hex: "#1877f2"),
            description: "Connect your Facebook account to access social features and posts",
            authorizationEndpoint: "https://www.facebook.com/v18.0/dialog/oauth",
            callbackPath: "/oauth2/facebook/callback",
            scopes: ["public_profile"],
            scopeSeparator: .comma,
            responseType: "code",
            isEnabled: true,
            additionalParameters: [:],
            requiresPKCE: false
        ),
        OAuthProvider(
            id: "spotify",
            name: "Spotify",
            iconName: "spotify",
            color: Color(hex: "#1db954"),
            description: "Connect your Spotify account to control music and playlists",
            authorizationEndpoint: "https://accounts.spotify.com/authorize",
            callbackPath: "/oauth2/spotify/callback",
            scopes: ["user-read-email", "user-read-private"],
            scopeSeparator: .space,
            responseType: "code",
            isEnabled: true,
            additionalParameters: [
                "show_dialog": "true"
            ],
            requiresPKCE: false
        ),
        OAuthProvider(
            id: "twitter",
            name: "Twitter / X",
            iconName: "twitter",
            color: Color(hex: "#1DA1F2"),
            description: "Connect your Twitter/X account to post tweets and access your timeline",
            authorizationEndpoint: "https://twitter.com/i/oauth2/authorize",
            callbackPath: "/oauth2/twitter/callback",
            scopes: ["tweet.read", "tweet.write", "users.read", "offline.access"],
            scopeSeparator: .space,
            responseType: "code",
            isEnabled: true,
            additionalParameters: [:],
            requiresPKCE: true
        )
    ]
}


