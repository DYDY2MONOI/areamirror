//
//  ContentView.swift
//  pocswift
//
//  Created by Dydy2Brazil on 18/09/2025.
//

import SwiftUI

struct ContentView: View {
    @State private var scale: CGFloat = 1.0

    var body: some View {
        VStack(spacing: 30) {
            Text("SwiftUI POC")
                .font(.largeTitle)

            Circle()
                .fill(Color.blue)
                .frame(width: 100, height: 100)
                .scaleEffect(scale)
                .onTapGesture {
                    withAnimation(.spring()) {
                        scale = scale == 1.0 ? 1.5 : 1.0
                    }
                }
        }
    }
}

#Preview {
    ContentView()
}
