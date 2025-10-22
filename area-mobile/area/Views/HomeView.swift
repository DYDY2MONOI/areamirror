//
//  HomeView.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct HomeView: View {
    @StateObject private var areaService = AreaService.shared
    @State private var showTestView = false
    @State private var showNewArea = false
    @State private var selectedTemplate: AreaTemplate?
    @State private var selectedArea: Area?
    @State private var selectedTab = 0
    let onLogout: () -> Void

    init(onLogout: @escaping () -> Void) {
        self.onLogout = onLogout
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
                                ProfileAvatar(size: 32, user: AuthService.shared.currentUser)

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
                            if areaService.isLoading {
                                ProgressView("Loading areas...")
                                    .foregroundColor(.white)
                                    .padding()
                            } else {
                                if areaService.userAreasLoaded {
                                    if !areaService.popularAreas.isEmpty {
                                        AppletSection(
                                            title: "Popular AREAs",
                                            applets: areaService.popularAreas.map { convertAreaTemplateToApplet($0) }
                                        )
                                    }

                                    if !areaService.recommendedAreas.isEmpty {
                                        AppletSection(
                                            title: "Recommended for you",
                                            applets: areaService.recommendedAreas.map { convertAreaTemplateToApplet($0) }
                                        )
                                    }
                                }

                                if selectedTab == 1 && !areaService.userAreas.isEmpty {
                                    AppletSection(
                                        title: "My AREAs",
                                        applets: areaService.userAreas.map { convertAreaToApplet($0) }
                                    )
                                }
                            }

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
        .fullScreenCover(item: $selectedTemplate) { template in
            EditAreaView(template: template)
        }
        .fullScreenCover(item: $selectedArea) { area in
            EditAreaView(area: area)
        }
        .onAppear {
            Task {
                await areaService.fetchAllAreas()
            }
        }
    }

    private func convertAreaToApplet(_ area: Area) -> Applet {
        return Applet(
            title: area.name,
            subtitle: "\(area.triggerService) → \(area.actionService)",
            description: area.description,
            icon: getServiceIcon(area.triggerService),
            gradient: getServiceGradient(area.triggerService, area.actionService),
            type: .create,
            action: { selectedArea = area }
        )
    }

    private func convertAreaTemplateToApplet(_ template: AreaTemplate) -> Applet {
        return Applet(
            title: template.title,
            subtitle: template.subtitle,
            description: template.description,
            icon: getServiceIcon(template.triggerService),
            gradient: getServiceGradient(template.triggerService, template.actionService),
            type: .create,
            action: {
                print("🔍 Checking for existing area for template: \(template.title)")
                print("🔍 User areas loaded: \(areaService.userAreasLoaded)")
                print("🔍 User areas count: \(areaService.userAreas.count)")
                print("🔍 Template trigger: \(template.triggerService), action: \(template.actionService)")
                
                for (index, area) in areaService.userAreas.enumerated() {
                    print("🔍 User area \(index): \(area.name) - \(area.triggerService) -> \(area.actionService)")
                }
                
                if let existingArea = areaService.userAreas.first(where: { 
                    $0.triggerService == template.triggerService && 
                    $0.actionService == template.actionService 
                }) {
                    print("🔄 Found existing area for template: \(template.title) -> Opening for edit")
                    print("🔄 Selected area ID: \(existingArea.id)")
                    selectedArea = existingArea
                } else {w area
                    print("➕ No existing area found for template: \(template.title) -> Opening for create")
                    selectedTemplate = template
                }
            }
        )
    }

    private func getServiceIcon(_ service: String) -> String {
        switch service.lowercased() {
        case "github": return "hammer.fill"
        case "gmail": return "envelope.fill"
        case "discord": return "message.fill"
        case "slack": return "message.circle.fill"
        case "weather": return "cloud.sun.fill"
        case "instagram": return "camera.fill"
        case "twitter": return "bird.fill"
        case "youtube": return "play.rectangle.fill"
        case "spotify": return "music.note"
        case "telegram": return "paperplane.fill"
        case "dropbox": return "folder.fill"
        case "notion": return "doc.text.fill"
        default: return "gear.fill"
        }
    }

    private func getServiceGradient(_ trigger: String, _ action: String) -> LinearGradient {
        let triggerColor = getServiceColor(trigger)
        let actionColor = getServiceColor(action)

        return LinearGradient(
            gradient: Gradient(colors: [triggerColor, actionColor]),
            startPoint: .topLeading,
            endPoint: .bottomTrailing
        )
    }

    private func getServiceColor(_ service: String) -> Color {
        switch service.lowercased() {
        case "github": return Color(red: 0.2, green: 0.2, blue: 0.2)
        case "gmail": return Color(red: 0.92, green: 0.26, blue: 0.21)
        case "discord": return Color(red: 0.35, green: 0.4, blue: 0.95)
        case "slack": return Color(red: 0.36, green: 0.8, blue: 0.36)
        case "weather": return Color(red: 0.0, green: 0.7, blue: 1.0)
        case "instagram": return Color(red: 0.8, green: 0.2, blue: 0.6)
        case "twitter": return Color(red: 0.1, green: 0.6, blue: 0.9)
        case "youtube": return Color(red: 1.0, green: 0.0, blue: 0.0)
        case "spotify": return Color(red: 0.2, green: 0.8, blue: 0.2)
        case "telegram": return Color(red: 0.0, green: 0.7, blue: 0.9)
        case "dropbox": return Color(red: 0.0, green: 0.5, blue: 0.8)
        case "notion": return Color(red: 0.2, green: 0.2, blue: 0.2)
        default: return Color.blue
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