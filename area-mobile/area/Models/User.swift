//
//  User.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation

struct User: Codable, Identifiable {
    let id: Int
    let email: String
    let firstName: String?
    let lastName: String?
    let createdAt: String?
    let updatedAt: String?
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
    let spotifyID: String?
    let spotifyEmail: String?
    let twitterID: String?
    let twitterUsername: String?
    
    enum CodingKeys: String, CodingKey {
        case id
        case email
        case firstName = "first_name"
        case lastName = "last_name"
        case createdAt = "created_at"
        case updatedAt = "updated_at"
        case phone
        case birthday
        case gender
        case country
        case lang
        case loginProvider = "login_provider"
        case profileImage = "profile_image"
        case githubID = "github_id"
        case githubUsername = "github_username"
        case googleID = "google_id"
        case googleEmail = "google_email"
        case facebookID = "facebook_id"
        case facebookEmail = "facebook_email"
        case spotifyID = "spotify_id"
        case spotifyEmail = "spotify_email"
        case twitterID = "twitter_id"
        case twitterUsername = "twitter_username"
    }
}

struct AuthResponse: Codable {
    let message: String
    let token: String
    let user: User
}

struct OAuthTokenResponse: Codable {
    let accessToken: String
    let refreshToken: String?
    let tokenType: String
    let expiresIn: Int
    let user: User

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case refreshToken = "refresh_token"
        case tokenType = "token_type"
        case expiresIn = "expires_in"
        case user
    }
}

struct RefreshTokenRequest: Codable {
    let refreshToken: String

    enum CodingKeys: String, CodingKey {
        case refreshToken = "refresh_token"
    }
}

struct RefreshTokenResponse: Codable {
    let accessToken: String
    let tokenType: String
    let expiresIn: Int

    enum CodingKeys: String, CodingKey {
        case accessToken = "access_token"
        case tokenType = "token_type"
        case expiresIn = "expires_in"
    }
}

struct LoginRequest: Codable {
    let email: String
    let password: String
}

struct RegisterRequest: Codable {
    let email: String
    let password: String
    let firstName: String?
    let lastName: String?
    
    enum CodingKeys: String, CodingKey {
        case email
        case password
        case firstName = "first_name"
        case lastName = "last_name"
    }
}

struct ProfileUpdateRequest: Codable {
    let firstName: String?
    let lastName: String?
    let phone: String?
    let country: String?
    let currentPassword: String?
    let newPassword: String?
    
    enum CodingKeys: String, CodingKey {
        case firstName = "first_name"
        case lastName = "last_name"
        case phone
        case country
        case currentPassword = "current_password"
        case newPassword = "new_password"
    }
}

struct ProfileResponse: Codable {
    let user: User
}
