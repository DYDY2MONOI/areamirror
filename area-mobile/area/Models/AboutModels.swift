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
    let current_time: Int
    let services: [AboutService]
}

struct AboutService: Decodable, Identifiable, Hashable {
    var id: String { name }
    let name: String
    let actions: [AboutAction]
    let reactions: [AboutReaction]
}

struct AboutAction: Decodable, Identifiable, Hashable {
    var id: String { name }
    let name: String
    let description: String
}

struct AboutReaction: Decodable, Identifiable, Hashable {
    var id: String { name }
    let name: String
    let description: String
}


