//
//  Color+Hex.swift
//  area
//
//  Created by GPT-5 Codex on 28/10/2025.
//

import SwiftUI

extension Color {
    init(hex: String) {
        let cleanedHex = hex.trimmingCharacters(in: CharacterSet.whitespacesAndNewlines.union(CharacterSet(charactersIn: "#"))).lowercased()

        var int: UInt64 = 0
        Scanner(string: cleanedHex).scanHexInt64(&int)

        let r, g, b, a: UInt64
        switch cleanedHex.count {
        case 8:
            (r, g, b, a) = (
                (int & 0xff000000) >> 24,
                (int & 0x00ff0000) >> 16,
                (int & 0x0000ff00) >> 8,
                int & 0x000000ff
            )
        case 6:
            (r, g, b, a) = (
                (int & 0xff0000) >> 16,
                (int & 0x00ff00) >> 8,
                int & 0x0000ff,
                0xff
            )
        default:
            (r, g, b, a) = (0, 0, 0, 0xff)
        }

        self.init(
            .sRGB,
            red: Double(r) / 255.0,
            green: Double(g) / 255.0,
            blue: Double(b) / 255.0,
            opacity: Double(a) / 255.0
        )
    }
}


