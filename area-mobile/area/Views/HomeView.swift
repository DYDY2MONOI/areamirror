//
//  HomeView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct HomeView: View {
    @State private var showTestView = false
    @State private var showNewArea = false
    @State private var selectedTab = 0
    let onLogout: () -> Void
    
    private let popularAreas = Applet.testApplets
    private let recommendedAreas = Applet.testApplets
    
    init(onLogout: @escaping () -> Void) {
        self.onLogout = onLogout
        print("HomeView init - popularAreas count: \(popularAreas.count)")
        print("HomeView init - recommendedAreas count: \(recommendedAreas.count)")
    }
    
    var body: some View {
        GeometryReader { geometry in
            ZStack {
                Color.black
                    .ignoresSafeArea()
                
                ScrollView {
                    VStack(spacing: 0) {
                        VStack(spacing: 20) {
                            HStack {
                                Circle()
                                    .fill(OptimizedGradients.primaryGradient)
                                    .frame(width: 32, height: 32)
                                    .overlay(
                                        Image(systemName: "person.fill")
                                            .font(.system(size: 16))
                                            .foregroundColor(.white)
                                    )
                                
                                Spacer()
                                
                                Button(action: { showTestView = true }) {
                                    Image(systemName: "testtube.2")
                                        .font(.system(size: 20))
                                        .foregroundColor(.white)
                                }
                            }
                            .padding(.horizontal, 20)
                            .padding(.top, 10)
                            
                            HStack(spacing: 0) {
                                TabButton(title: "All", isSelected: selectedTab == 0) {
                                    selectedTab = 0
                                }
                                TabButton(title: "My AREAs", isSelected: selectedTab == 1) {
                                    selectedTab = 1
                                }
                                TabButton(title: "Popular", isSelected: selectedTab == 2) {
                                    selectedTab = 2
                                }
                                TabButton(title: "Create", isSelected: selectedTab == 3) {
                                    selectedTab = 3
                                }
                            }
                            .padding(.horizontal, 20)
                        }
                        
                        VStack(spacing: 32) {
                            AppletSection(
                                title: "Popular AREAs",
                                applets: popularAreas
                            )
                            
                            AppletSection(
                                title: "Recommended for you",
                                applets: recommendedAreas
                            )
                            
                            AppletSection(
                                title: "Create new AREA",
                                applets: [
                                    Applet(
                                        title: "New AREA",
                                        subtitle: "Get started",
                                        description: "Connect your favorite services",
                                        icon: "plus.circle.fill",
                                        gradient: OptimizedGradients.primaryGradient,
                                        type: .create,
                                        action: { showNewArea = true }
                                    ),
                                    Applet(
                                        title: "Email Template",
                                        subtitle: "Gmail automation",
                                        description: "Automate your important emails",
                                        icon: "envelope.badge.fill",
                                        gradient: OptimizedGradients.blueGradient,
                                        type: .create,
                                        action: { print("Email Template") }
                                    ),
                                    Applet(
                                        title: "Social Template",
                                        subtitle: "Social networks",
                                        description: "Automate your posts and shares",
                                        icon: "share.and.arrow.up.fill",
                                        gradient: OptimizedGradients.purpleGradient,
                                        type: .create,
                                        action: { print("Social Template") }
                                    )
                                ]
                            )
                        }
                        .padding(.top, 20)
                        .padding(.bottom, 40)
                    }
                }
            }
        }
        .fullScreenCover(isPresented: $showTestView) {
            TestView()
        }
        .fullScreenCover(isPresented: $showNewArea) {
            NewAreaView()
        }
    }
}

struct TabButton: View {
    let title: String
    let isSelected: Bool
    let action: () -> Void
    
    var body: some View {
        Button(action: action) {
            Text(title)
                .font(.system(size: 16, weight: .medium))
                .foregroundColor(isSelected ? .white : .gray)
                .padding(.horizontal, 16)
                .padding(.vertical, 8)
                .background(
                    RoundedRectangle(cornerRadius: 20)
                        .fill(isSelected ? Color.green : Color.clear)
                )
        }
    }
}

#Preview {
    HomeView(onLogout: {})
}
