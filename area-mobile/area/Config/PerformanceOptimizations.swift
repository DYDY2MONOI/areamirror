//
//  PerformanceOptimizations.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct OptimizedGradients {
    private static let gradientCache = NSCache<NSString, CAGradientLayer>()
    
    static let background = LinearGradient(
        gradient: Gradient(colors: [
            Color.black,
            AppColors.darkerBackground,
            Color(red: 0.1, green: 0.1, blue: 0.2)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let primaryGradient = LinearGradient(
        gradient: Gradient(colors: [
            AppColors.primaryBlue,
            AppColors.secondaryPurple
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let blueGradient = LinearGradient(
        gradient: Gradient(colors: [
            Color.blue.opacity(0.8),
            Color.cyan.opacity(0.6)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let purpleGradient = LinearGradient(
        gradient: Gradient(colors: [
            Color.purple.opacity(0.8),
            Color.pink.opacity(0.6)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let greenGradient = LinearGradient(
        gradient: Gradient(colors: [
            Color.green.opacity(0.8),
            Color.mint.opacity(0.6)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
    
    static let orangeGradient = LinearGradient(
        gradient: Gradient(colors: [
            Color.orange.opacity(0.8),
            Color.yellow.opacity(0.6)
        ]),
        startPoint: .topLeading,
        endPoint: .bottomTrailing
    )
}

struct OptimizedDecorativeShapes: View {
    var body: some View {
        ZStack {
            Circle()
                .fill(AppColors.primaryBlue.opacity(0.1))
                .frame(width: 200, height: 200)
                .offset(x: 100, y: -100)
            
            Circle()
                .fill(AppColors.secondaryPurple.opacity(0.1))
                .frame(width: 150, height: 150)
                .offset(x: -80, y: 120)
        }
        .drawingGroup()
    }
}

class PerformanceMonitor: ObservableObject {
    @Published var isHighPerformanceMode = false
    
    func enableHighPerformanceMode() {
        isHighPerformanceMode = true
    }
    
    func disableHighPerformanceMode() {
        isHighPerformanceMode = false
    }
}

struct OptimizedBackground: ViewModifier {
    func body(content: Content) -> some View {
        content
            .background(OptimizedGradients.background)
            .drawingGroup()
    }
}

struct OptimizedCard: ViewModifier {
    let gradient: LinearGradient
    
    func body(content: Content) -> some View {
        content
            .background(gradient)
            .cornerRadius(12)
            .drawingGroup()
    }
}

extension View {
    func optimizedBackground() -> some View {
        modifier(OptimizedBackground())
    }
    
    func optimizedCard(gradient: LinearGradient) -> some View {
        modifier(OptimizedCard(gradient: gradient))
    }
    
    func performanceOptimized() -> some View {
        self
            .drawingGroup()
            .clipped()
    }
}
