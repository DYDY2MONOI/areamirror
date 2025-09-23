//
//  SearchView.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI

struct SearchView: View {
    @State private var searchText = ""
    @State private var selectedCategory = "All"
    
    let categories = ["All", "Actions", "Reactions", "Services", "Templates"]
    
    var body: some View {
        NavigationView {
            ZStack {
                AppGradients.background
                    .ignoresSafeArea()
                
                VStack(spacing: 20) {
                    HStack {
                        Image(systemName: "magnifyingglass")
                            .foregroundColor(.gray)
                        
                        TextField("Search areas, services...", text: $searchText)
                            .textFieldStyle(PlainTextFieldStyle())
                            .autocapitalization(.none)
                            .disableAutocorrection(true)
                            .foregroundColor(.white)
                            .onSubmit {
                                performSearch()
                            }
                        
                        if !searchText.isEmpty {
                            Button("Clear") {
                                searchText = ""
                            }
                            .foregroundColor(AppColors.primaryBlue)
                        }
                    }
                    .padding(.horizontal, 16)
                    .padding(.vertical, 12)
                    .background(
                        RoundedRectangle(cornerRadius: 12)
                            .fill(AppColors.darkBackground)
                            .overlay(
                                RoundedRectangle(cornerRadius: 12)
                                    .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                            )
                    )
                    .padding(.horizontal, 20)
                    
                    ScrollView(.horizontal, showsIndicators: false) {
                        HStack(spacing: 12) {
                            ForEach(categories, id: \.self) { category in
                                Button(action: {
                                    selectedCategory = category
                                }) {
                                    Text(category)
                                        .font(.system(size: 14, weight: .medium))
                                        .foregroundColor(selectedCategory == category ? .white : .gray)
                                        .padding(.horizontal, 16)
                                        .padding(.vertical, 8)
                                        .background(
                                            RoundedRectangle(cornerRadius: 20)
                                                .fill(selectedCategory == category ? AppColors.primaryBlue : Color.clear)
                                                .overlay(
                                                    RoundedRectangle(cornerRadius: 20)
                                                        .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                                                )
                                        )
                                }
                            }
                        }
                        .padding(.horizontal, 20)
                    }
                    
                    if searchText.isEmpty {
                        VStack(spacing: 20) {
                            Image(systemName: "magnifyingglass.circle")
                                .font(.system(size: 60))
                                .foregroundColor(.gray)
                            
                            Text("Search Areas")
                                .font(.title2)
                                .fontWeight(.bold)
                                .foregroundColor(.white)
                            
                            Text("Find actions, reactions, and services to create your automation areas")
                                .font(.body)
                                .foregroundColor(.gray)
                                .multilineTextAlignment(.center)
                                .padding(.horizontal, 40)
                        }
                        .padding(.top, 60)
                    } else {
                        ScrollView {
                            LazyVStack(spacing: 16) {
                                ForEach(0..<10) { index in
                                    SearchResultCard(
                                        title: "Sample Area \(index + 1)",
                                        description: "This is a sample automation area",
                                        category: categories.randomElement() ?? "All"
                                    )
                                }
                            }
                            .padding(.horizontal, 20)
                        }
                    }
                    
                    Spacer()
                }
            }
        }
        .navigationTitle("Search")
        .navigationBarHidden(true)
    }
    
    private func performSearch() {
        print("Searching for: \(searchText) in category: \(selectedCategory)")
    }
}

struct SearchResultCard: View {
    let title: String
    let description: String
    let category: String
    
    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 8) {
                Text(title)
                    .font(.headline)
                    .foregroundColor(.white)
                
                Text(description)
                    .font(.caption)
                    .foregroundColor(.gray)
                
                Text(category)
                    .font(.caption2)
                    .foregroundColor(AppColors.primaryBlue)
                    .padding(.horizontal, 8)
                    .padding(.vertical, 4)
                    .background(
                        RoundedRectangle(cornerRadius: 8)
                            .fill(AppColors.primaryBlue.opacity(0.2))
                    )
            }
            
            Spacer()
            
            Button(action: {
            }) {
                Image(systemName: "plus.circle")
                    .font(.title2)
                    .foregroundColor(AppColors.primaryBlue)
            }
        }
        .padding(16)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(AppColors.darkBackground)
                .overlay(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(Color.gray.opacity(0.2), lineWidth: 1)
                )
        )
    }
}

#Preview {
    SearchView()
}
