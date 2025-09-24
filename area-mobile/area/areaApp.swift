//
//  areaApp.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

@main
struct areaApp: App {
    @StateObject private var performanceManager = PerformanceManager()
    
    var body: some Scene {
        WindowGroup {
            PerformanceOptimizedView {
                ContentView()
            }
            .environmentObject(performanceManager)
            .onAppear {
                if AppPerformanceConfig.enableHighPerformanceMode {
                    performanceManager.enableHighPerformanceMode()
                }
            }
        }
    }
}
