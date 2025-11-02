//
//  AppNotification.swift
//  area
//
//  Created by Codex on 02/11/2025.
//

import Foundation

struct AppNotification: Codable, Identifiable, Equatable {
    let id: UUID
    var title: String
    var body: String
    var date: Date
    var isRead: Bool
    var source: String?

    init(id: UUID = UUID(), title: String, body: String, date: Date = Date(), isRead: Bool = false, source: String? = nil) {
        self.id = id
        self.title = title
        self.body = body
        self.date = date
        self.isRead = isRead
        self.source = source
    }
}

