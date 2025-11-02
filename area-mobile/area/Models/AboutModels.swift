//
//  AboutModels.swift
//  area
//
//  Created to mirror backend /about.json for services/actions/reactions
//

import Foundation

struct AboutResponse: Decodable {
    let server: AboutServer
}

struct AboutServer: Decodable {
    let currentTime: Int
    let services: [AboutService]
    
    enum CodingKeys: String, CodingKey {
        case currentTime = "current_time"
        case services
    }
}

struct AboutService: Decodable, Identifiable, Hashable {
    var id: String { name }
    let name: String
    let actions: [AboutAction]
    let reactions: [AboutReaction]
}

struct AboutAction: Decodable, Identifiable, Hashable {
    var id: String { identifier ?? name }
    let name: String
    let description: String
    let identifier: String?
    
    // Initialize identifier as nil since backend doesn't provide it
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        name = try container.decode(String.self, forKey: .name)
        description = try container.decode(String.self, forKey: .description)
        identifier = nil // Backend doesn't send identifier
    }
    
    enum CodingKeys: String, CodingKey {
        case name
        case description
    }
}

struct AboutReaction: Decodable, Identifiable, Hashable {
    var id: String { identifier ?? name }
    let name: String
    let description: String
    let identifier: String?
    
    // Initialize identifier as nil since backend doesn't provide it
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        name = try container.decode(String.self, forKey: .name)
        description = try container.decode(String.self, forKey: .description)
        identifier = nil // Backend doesn't send identifier
    }
    
    enum CodingKeys: String, CodingKey {
        case name
        case description
    }
}


