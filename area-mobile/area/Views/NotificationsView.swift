//
//  NotificationsView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import SwiftUI

struct NotificationsView: View {
    @Environment(\.dismiss) private var dismiss
    @ObservedObject private var notificationService = NotificationService.shared

    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()

                if notificationService.notifications.isEmpty {
                    VStack(spacing: 12) {
                        Image(systemName: "bell.slash")
                            .font(.system(size: 48))
                            .foregroundColor(.gray)
                        Text("No notifications yet")
                            .foregroundColor(.gray)
                            .font(.headline)
                        Text("You’ll see AREA alerts here.")
                            .foregroundColor(.gray)
                            .font(.subheadline)
                    }
                } else {
                    List {
                        ForEach(notificationService.notifications) { item in
                            VStack(alignment: .leading, spacing: 6) {
                                HStack {
                                    Text(item.title)
                                        .font(.headline)
                                        .foregroundColor(.white)
                                    if !item.isRead {
                                        Circle()
                                            .fill(AppColors.primaryBlue)
                                            .frame(width: 8, height: 8)
                                    }
                                }
                                Text(item.body)
                                    .font(.subheadline)
                                    .foregroundColor(.gray)
                                HStack {
                                    if let src = item.source, !src.isEmpty {
                                        Text(src)
                                            .font(.caption)
                                            .foregroundColor(.gray)
                                    }
                                    Spacer()
                                    Text(item.date.formatted(date: .abbreviated, time: .shortened))
                                        .font(.caption)
                                        .foregroundColor(.gray)
                                }
                            }
                            .listRowBackground(AppColors.darkBackground)
                            .padding(.vertical, 6)
                        }
                        .onDelete { offsets in
                            notificationService.remove(at: offsets)
                        }
                    }
                    .listStyle(.plain)
                    .scrollContentBackground(.hidden)
                }
            }
            .navigationTitle("Notifications")
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button(action: { dismiss() }) {
                        Image(systemName: "xmark")
                    }
                    .foregroundColor(.white)
                }
                ToolbarItem(placement: .navigationBarTrailing) {
                    Menu {
                        Button("Mark all as read") { notificationService.markAllAsRead() }
                        Button("Clear all", role: .destructive) { notificationService.removeAll() }
                    } label: {
                        Image(systemName: "ellipsis.circle")
                    }
                    .foregroundColor(.white)
                }
            }
        }
        .navigationViewStyle(.stack)
        .task {
            await NotificationService.shared.importDeliveredNotifications()
        }
    }
}

#Preview {
    NotificationsView()
}
