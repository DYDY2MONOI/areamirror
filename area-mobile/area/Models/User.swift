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
    }
}

struct AuthResponse: Codable {
    let message: String
    let token: String
    let user: User
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
