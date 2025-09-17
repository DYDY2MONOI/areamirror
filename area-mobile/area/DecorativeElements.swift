//
//  DecorativeElements.swift
//  area
//
//  Created by Dydy2Brazil on 16/09/2025.
//

import SwiftUI

struct DecorativeShapes: View {
    var body: some View {
        VStack {
            HStack {
                Spacer()
                RoundedRectangle(cornerRadius: 50)
                    .fill(AppGradients.decorativeBlue)
                    .frame(width: 120, height: 80)
                    .rotationEffect(.degrees(15))
                    .offset(x: 20, y: -50)
            }
            Spacer()
            HStack {
                RoundedRectangle(cornerRadius: 40)
                    .fill(AppGradients.decorativePurple)
                    .frame(width: 100, height: 60)
                    .rotationEffect(.degrees(-20))
                    .offset(x: -30, y: 50)
                Spacer()
            }
        }
    }
}
