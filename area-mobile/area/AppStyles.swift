//
//  AppStyles.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct AppColors {
    static let primaryBlue = Color(red: 0.34, green: 0.50, blue: 0.91)
    static let secondaryPurple = Color(red: 0.53, green: 0.38, blue: 0.82)
    static let lightBlue = Color(red: 0.52, green: 0.81, blue: 0.92)
    static let lightPurple = Color(red: 0.76, green: 0.78, blue: 0.89)
    static let darkBackground = Color(red: 0.1, green: 0.1, blue: 0.15)
    static let darkerBackground = Color(red: 0.05, green: 0.05, blue: 0.15)
}

struct AppGradients {
    static let background = LinearGradient(
        gradient: Gradient(colors: [
            Color.black,
            AppColors.darkerBackground,
            Color(red: 0.1, green: 0.1, blue: 0.2)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let button = LinearGradient(
        gradient: Gradient(colors: [
            AppColors.primaryBlue,
            AppColors.secondaryPurple
        ]),
        startPoint: .leading,
        endPoint: .trailing
    )
    
    static let decorativeBlue = LinearGradient(
        gradient: Gradient(colors: [
            AppColors.primaryBlue.opacity(0.3),
            AppColors.lightBlue.opacity(0.2)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let decorativePurple = LinearGradient(
        gradient: Gradient(colors: [
            AppColors.secondaryPurple.opacity(0.3),
            AppColors.lightPurple.opacity(0.2)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
}

struct AppTextStyles {
    static let title = Font.system(size: 36, weight: .bold, design: .rounded)
    static let subtitle = Font.system(size: 24, weight: .bold)
    static let body = Font.system(size: 16)
    static let caption = Font.system(size: 14)
    static let button = Font.system(size: 16, weight: .bold)
}
