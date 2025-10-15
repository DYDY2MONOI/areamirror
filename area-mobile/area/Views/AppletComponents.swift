//
//  AppletComponents.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct AppletCard: View {
    let applet: Applet
    
    var body: some View {
        Button(action: applet.action) {
            VStack(alignment: .leading, spacing: 0) {
                ZStack {
                    RoundedRectangle(cornerRadius: 12)
                        .fill(applet.gradient)
                        .frame(width: 160, height: 160)
                    
                    VStack {
                        if applet.type == .create {
                            Image(systemName: applet.icon)
                                .font(.system(size: 40, weight: .bold))
                                .foregroundColor(.white)
                        } else if applet.type == .blend {
                            HStack(spacing: -10) {
                                Circle()
                                    .fill(Color.orange)
                                    .frame(width: 30, height: 30)
                                Circle()
                                    .fill(Color.yellow)
                                    .frame(width: 30, height: 30)
                            }
                            .overlay(
                                Image(systemName: "music.note")
                                    .font(.system(size: 12, weight: .bold))
                                    .foregroundColor(.white)
                                    .offset(x: 15, y: -15)
                            )
                        } else {
                            Image(systemName: applet.icon)
                                .font(.system(size: 30, weight: .bold))
                                .foregroundColor(.white)
                        }
                    }
                }
                
                VStack(alignment: .leading, spacing: 4) {
                    Text(applet.title)
                        .font(.system(size: 14, weight: .bold))
                        .foregroundColor(.white)
                        .lineLimit(1)
                    
                    Text(applet.subtitle)
                        .font(.system(size: 12, weight: .medium))
                        .foregroundColor(.gray)
                        .lineLimit(1)
                    
                    Text(applet.description)
                        .font(.system(size: 11))
                        .foregroundColor(.gray)
                        .lineLimit(2)
                }
                .padding(.top, 8)
                .frame(width: 160, alignment: .leading)
            }
        }
        .buttonStyle(PlainButtonStyle())
        .performanceOptimized()
    }
}

struct AppletSection: View {
    let title: String
    let applets: [Applet]
    
    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            Text(title)
                .font(.system(size: 22, weight: .bold))
                .foregroundColor(.white)
                .padding(.horizontal, 20)
                .onAppear {
                    print("AppletSection '\(title)' rendering with \(applets.count) applets")
                }
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 16) {
                    ForEach(applets) { applet in
                        AppletCard(applet: applet)
                    }
                }
                .padding(.horizontal, 20)
            }
        }
        .performanceOptimized()
    }
}

struct NowPlayingBar: View {
    var body: some View {
        HStack(spacing: 12) {
            RoundedRectangle(cornerRadius: 4)
                .fill(Color.gray.opacity(0.3))
                .frame(width: 50, height: 50)
                .overlay(
                    Image(systemName: "music.note")
                        .font(.title2)
                        .foregroundColor(.white)
                )
            
            VStack(alignment: .leading, spacing: 2) {
                Text("Il Est Le Roi")
                    .font(.system(size: 14, weight: .medium))
                    .foregroundColor(.white)
                    .lineLimit(1)
                
                HStack(spacing: 4) {
                    Text("Fortuné")
                        .font(.system(size: 12))
                        .foregroundColor(.gray)
                    
                    Image(systemName: "diamond.fill")
                        .font(.system(size: 8))
                        .foregroundColor(.green)
                }
            }
            
            Spacer()
            
            HStack(spacing: 16) {
                Button(action: {}) {
                    Image(systemName: "speaker.wave.2")
                        .font(.system(size: 16))
                        .foregroundColor(.gray)
                }
                
                Button(action: {}) {
                    Image(systemName: "play.fill")
                        .font(.system(size: 16))
                        .foregroundColor(.white)
                }
            }
        }
        .padding(.horizontal, 16)
        .padding(.vertical, 12)
        .background(
            RoundedRectangle(cornerRadius: 8)
                .fill(Color.red.opacity(0.9))
        )
        .padding(.horizontal, 20)
    }
}

#Preview(traits: .sizeThatFitsLayout) {
    VStack {
        AppletCard(applet: Applet.sampleApplets[0])
        AppletCard(applet: Applet.sampleApplets[1])
        AppletCard(applet: Applet.sampleApplets[3])
    }
    .background(Color.black)
}
