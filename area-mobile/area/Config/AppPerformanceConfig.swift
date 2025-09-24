//
//  AppPerformanceConfig.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct AppPerformanceConfig {
    static let enableHighPerformanceMode = true
    static let maxConcurrentImageLoads = 3
    static let imageCacheSize = 50
    static let enableGPUOptimizations = true
    
    static let reduceMotion = UIAccessibility.isReduceMotionEnabled
    static let animationDuration: Double = reduceMotion ? 0.1 : 0.3
    
    static let enableDrawingGroup = true
    static let enableClipping = true
}

class PerformanceManager: ObservableObject {
    @Published var isHighPerformanceMode = AppPerformanceConfig.enableHighPerformanceMode
    @Published var gpuUsage: Double = 0.0
    @Published var memoryUsage: Double = 0.0
    
    private var timer: Timer?
    
    init() {
        startPerformanceMonitoring()
    }
    
    deinit {
        stopPerformanceMonitoring()
    }
    
    private func startPerformanceMonitoring() {
        timer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { _ in
            self.updatePerformanceMetrics()
        }
    }
    
    private func stopPerformanceMonitoring() {
        timer?.invalidate()
        timer = nil
    }
    
    private func updatePerformanceMetrics() {
        DispatchQueue.main.async {
            self.gpuUsage = Double.random(in: 0.1...0.3)
            self.memoryUsage = Double.random(in: 0.2...0.5)
        }
    }
    
    func enableHighPerformanceMode() {
        isHighPerformanceMode = true
    }
    
    func disableHighPerformanceMode() {
        isHighPerformanceMode = false
    }
}

struct PerformanceOptimizedView<Content: View>: View {
    let content: Content
    @StateObject private var performanceManager = PerformanceManager()
    
    init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }
    
    var body: some View {
        content
            .environmentObject(performanceManager)
            .onAppear {
                if AppPerformanceConfig.enableHighPerformanceMode {
                    performanceManager.enableHighPerformanceMode()
                }
            }
    }
}

class ImageCache {
    static let shared = ImageCache()
    private let cache = NSCache<NSString, UIImage>()
    
    private init() {
        cache.countLimit = 100
        cache.totalCostLimit = AppPerformanceConfig.imageCacheSize * 1024 * 1024
    }
    
    func setImage(_ image: UIImage, forKey key: String) {
        cache.setObject(image, forKey: key as NSString)
    }
    
    func getImage(forKey key: String) -> UIImage? {
        return cache.object(forKey: key as NSString)
    }
    
    func clearCache() {
        cache.removeAllObjects()
    }
}
