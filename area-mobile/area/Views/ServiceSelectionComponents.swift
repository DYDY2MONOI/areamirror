//
//  ServiceSelectionComponents.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct ServiceCard: View {
    let service: Service
    let isSelected: Bool
    let onTap: () -> Void
    
    var body: some View {
        Button(action: onTap) {
            VStack(spacing: 8) {
                ZStack {
                    Circle()
                        .fill(service.color.opacity(0.2))
                        .frame(width: 50, height: 50)
                    
                    ZStack {
                        RoundedRectangle(cornerRadius: 4)
                            .fill(service.color == .white ? Color.gray : service.color)
                            .frame(width: 24, height: 24)
                        
                        Image(service.icon)
                            .resizable()
                            .aspectRatio(contentMode: .fit)
                            .frame(width: 20, height: 20)
                            .colorMultiply(service.color == .white ? .white : service.color)
                    }
                }
                
                Text(service.name)
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.white)
                    .multilineTextAlignment(.center)
                    .lineLimit(2)
            }
            .frame(width: 80, height: 90)
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(isSelected ? service.color.opacity(0.3) : Color.gray.opacity(0.1))
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(isSelected ? service.color : Color.clear, lineWidth: 2)
                    )
            )
        }
        .buttonStyle(PlainButtonStyle())
    }
}

struct ServiceCategorySection: View {
    let category: Service.ServiceCategory
    let services: [Service]
    @Binding var selectedService: Service?
    let onServiceSelected: (Service) -> Void
    
    var categoryTitle: String {
        switch category {
        case .social: return "Social Media"
        case .productivity: return "Productivity"
        case .communication: return "Communication"
        case .entertainment: return "Entertainment"
        case .development: return "Development"
        case .other: return "Other"
        }
    }
    
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text(categoryTitle)
                .font(.system(size: 18, weight: .bold))
                .foregroundColor(.white)
                .padding(.horizontal, 20)
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    ForEach(services) { service in
                        ServiceCard(
                            service: service,
                            isSelected: selectedService?.id == service.id,
                            onTap: {
                                onServiceSelected(service)
                            }
                        )
                    }
                }
                .padding(.horizontal, 20)
            }
        }
    }
}

struct ServiceSelectionView: View {
    let title: String
    @Binding var selectedService: Service?
    let onServiceSelected: (Service) -> Void
    
    @State private var selectedCategory: Service.ServiceCategory = .social
    
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text(title)
                .font(.system(size: 20, weight: .bold))
                .foregroundColor(.white)
                .padding(.horizontal, 20)
            
            ScrollView(.horizontal, showsIndicators: false) {
                HStack(spacing: 12) {
                    ForEach([Service.ServiceCategory.social, .productivity, .communication, .entertainment, .development, .other], id: \.self) { category in
                        Button(action: {
                            selectedCategory = category
                        }) {
                            Text(categoryTitle(for: category))
                                .font(.system(size: 14, weight: .medium))
                                .foregroundColor(selectedCategory == category ? .white : .gray)
                                .padding(.horizontal, 16)
                                .padding(.vertical, 8)
                                .background(
                                    RoundedRectangle(cornerRadius: 20)
                                        .fill(selectedCategory == category ? AppColors.primaryBlue : Color.gray.opacity(0.2))
                                )
                        }
                    }
                }
                .padding(.horizontal, 20)
            }
            
            ServiceCategorySection(
                category: selectedCategory,
                services: Service.services(for: selectedCategory),
                selectedService: $selectedService,
                onServiceSelected: onServiceSelected
            )
        }
    }
    
    private func categoryTitle(for category: Service.ServiceCategory) -> String {
        switch category {
        case .social: return "Social"
        case .productivity: return "Productivity"
        case .communication: return "Communication"
        case .entertainment: return "Entertainment"
        case .development: return "Development"
        case .other: return "Other"
        }
    }
}

#Preview {
    VStack {
        ServiceSelectionView(
            title: "Select Action Service",
            selectedService: .constant(nil),
            onServiceSelected: { _ in }
        )
    }
    .background(Color.black)
    .previewLayout(.sizeThatFits)
}
