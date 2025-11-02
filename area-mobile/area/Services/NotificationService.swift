//
//  NotificationService.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.≈
//

import Foundation
import UserNotifications
import UIKit

@MainActor
final class NotificationService: NSObject, ObservableObject {
    static let shared = NotificationService()

    @Published private(set) var notifications: [AppNotification] = []

    private let storageKey = "app_notifications_store_v1"
    private let deviceTokenKey = "apns_device_token_string"

    private override init() {
        super.init()
        loadSavedNotifications()
    }

    func configure() {
        let center = UNUserNotificationCenter.current()
        center.delegate = self

        center.requestAuthorization(options: [.alert, .badge, .sound]) { [weak self] granted, error in
            if let error = error {
                print("Notification permission error: \(error.localizedDescription)")
            }
            print("Notification permission granted: \(granted)")
            if granted {
                DispatchQueue.main.async {
                    UIApplication.shared.registerForRemoteNotifications()
                }
            }
            Task { @MainActor in
                self?.appendInternalNotice(title: granted ? "Notifications enabled" : "Notifications disabled",
                                            body: granted ? "You’ll receive alerts when areas trigger." : "Enable notifications in Settings to get alerts.")
                await self?.importDeliveredNotifications()
            }
        }
    }

    func didRegisterForRemoteNotifications(deviceToken: Data) {
        let tokenString = deviceToken.map { String(format: "%02.2hhx", $0) }.joined()
        UserDefaults.standard.set(tokenString, forKey: deviceTokenKey)
        print("APNs device token: \(tokenString)")

        Task { @MainActor in
            await registerDeviceTokenWithBackend(tokenString)
        }
    }

    func didFailToRegisterForRemoteNotifications(error: Error) {
        print("Failed to register for remote notifications: \(error.localizedDescription)")
    }

    func handleIncoming(userInfo: [AnyHashable: Any]) {
        var title = AppConfig.appName
        var body = "New notification"
        var source: String? = nil

        if let aps = userInfo["aps"] as? [String: Any] {
            if let alert = aps["alert"] as? [String: Any] {
                if let t = alert["title"] as? String, !t.isEmpty { title = t }
                if let b = alert["body"] as? String, !b.isEmpty { body = b }
            } else if let alertString = aps["alert"] as? String {
                body = alertString
            }
        }

        if let src = userInfo["source"] as? String { source = src }

        let item = AppNotification(title: title, body: body, source: source)
        notifications.insert(item, at: 0)
        persistNotifications()
    }

    func importDeliveredNotifications() async {
        let center = UNUserNotificationCenter.current()
        let delivered = await withCheckedContinuation { (cont: CheckedContinuation<[UNNotification], Never>) in
            center.getDeliveredNotifications { notes in
                cont.resume(returning: notes)
            }
        }
        guard !delivered.isEmpty else { return }
        for note in delivered {
            let c = note.request.content
            var title = c.title
            let body = c.body
            if title.isEmpty { title = AppConfig.appName }
            let src = c.userInfo["source"] as? String
            let item = AppNotification(title: title, body: body, source: src)
            notifications.insert(item, at: 0)
        }
        persistNotifications()
    }

    func appendInternalNotice(title: String, body: String) {
        let item = AppNotification(title: title, body: body, source: "system")
        notifications.insert(item, at: 0)
        persistNotifications()
    }

    func markAllAsRead() {
        notifications = notifications.map { n in
            var nn = n
            nn.isRead = true
            return nn
        }
        persistNotifications()
    }

    func removeAll() {
        notifications.removeAll()
        persistNotifications()
    }

    func remove(at offsets: IndexSet) {
        notifications.remove(atOffsets: offsets)
        persistNotifications()
    }

    private func persistNotifications() {
        do {
            let data = try JSONEncoder().encode(notifications)
            UserDefaults.standard.set(data, forKey: storageKey)
        } catch {
            print("Failed to persist notifications: \(error.localizedDescription)")
        }
    }

    private func loadSavedNotifications() {
        guard let data = UserDefaults.standard.data(forKey: storageKey) else { return }
        do {
            let saved = try JSONDecoder().decode([AppNotification].self, from: data)
            notifications = saved
        } catch {
            print("Failed to load saved notifications: \(error.localizedDescription)")
        }
    }

    private func registerDeviceTokenWithBackend(_ tokenString: String) async {
        guard let header = AuthService.shared.authorizationHeader() else { return }
        guard let url = URL(string: AppConfig.getAPIEndpoint("/mobile/push/register")) else { return }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(header, forHTTPHeaderField: "Authorization")

        let payload: [String: Any] = [
            "device_token": tokenString,
            "platform": "ios"
        ]

        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: payload)
            let (data, response) = try await URLSession.shared.data(for: request)
            if let http = response as? HTTPURLResponse {
                print("/mobile/push/register status: \(http.statusCode)")
                if !(200...299).contains(http.statusCode) {
                    let body = String(data: data, encoding: .utf8) ?? "<no-body>"
                    print("Failed to register device token: \(body)")
                }
            }
        } catch {
            print("Register device token request failed: \(error.localizedDescription)")
        }
    }
}

extension NotificationService: UNUserNotificationCenterDelegate {
    nonisolated func userNotificationCenter(_ center: UNUserNotificationCenter,
                                            willPresent notification: UNNotification,
                                            withCompletionHandler completionHandler: @escaping (UNNotificationPresentationOptions) -> Void) {
        completionHandler([.banner, .sound, .badge])

        let content = notification.request.content
        let info = notification.request.content.userInfo

        Task { @MainActor in
            var source: String? = nil
            if let src = info["source"] as? String { source = src }
            let item = AppNotification(title: content.title.isEmpty ? AppConfig.appName : content.title,
                                       body: content.body,
                                       source: source)
            self.notifications.insert(item, at: 0)
            self.persistNotifications()
        }
    }

    nonisolated func userNotificationCenter(_ center: UNUserNotificationCenter,
                                            didReceive response: UNNotificationResponse,
                                            withCompletionHandler completionHandler: @escaping () -> Void) {
        let content = response.notification.request.content
        let info = content.userInfo
        Task { @MainActor in
            var source: String? = nil
            if let src = info["source"] as? String { source = src }
            let item = AppNotification(title: content.title.isEmpty ? AppConfig.appName : content.title,
                                       body: content.body,
                                       isRead: true,
                                       source: source)
            self.notifications.insert(item, at: 0)
            self.persistNotifications()
        }
        completionHandler()
    }
}
