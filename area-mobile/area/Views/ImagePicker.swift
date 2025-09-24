//
//  ImagePicker.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import SwiftUI
import UIKit

struct ImagePicker: UIViewControllerRepresentable {
    @Binding var selectedImage: UIImage?
    @Environment(\.presentationMode) var presentationMode
    var sourceType: UIImagePickerController.SourceType = .photoLibrary
    
    func makeUIViewController(context: Context) -> UIImagePickerController {
        let picker = UIImagePickerController()
        picker.delegate = context.coordinator
        picker.sourceType = sourceType
        picker.allowsEditing = true
        return picker
    }
    
    func updateUIViewController(_ uiViewController: UIImagePickerController, context: Context) {}
    
    func makeCoordinator() -> Coordinator {
        Coordinator(self)
    }
    
    class Coordinator: NSObject, UIImagePickerControllerDelegate, UINavigationControllerDelegate {
        let parent: ImagePicker
        
        init(_ parent: ImagePicker) {
            self.parent = parent
        }
        
        func imagePickerController(_ picker: UIImagePickerController, didFinishPickingMediaWithInfo info: [UIImagePickerController.InfoKey : Any]) {
            if let editedImage = info[.editedImage] as? UIImage {
                parent.selectedImage = editedImage
            } else if let originalImage = info[.originalImage] as? UIImage {
                parent.selectedImage = originalImage
            }
            
            parent.presentationMode.wrappedValue.dismiss()
        }
        
        func imagePickerControllerDidCancel(_ picker: UIImagePickerController) {
            parent.presentationMode.wrappedValue.dismiss()
        }
    }
}

struct ProfileImagePicker: View {
    @Binding var selectedImage: UIImage?
    @State private var showingImagePicker = false
    @State private var showingActionSheet = false
    @State private var imagePickerSourceType: UIImagePickerController.SourceType = .photoLibrary
    @StateObject private var authService = AuthService.shared
    @State private var showingPermissionAlert = false
    @State private var permissionMessage = ""
    
    var body: some View {
        VStack(spacing: 16) {
            Button(action: {
                showingActionSheet = true
            }) {
                ZStack {
                    if let selectedImage = selectedImage {
                        Image(uiImage: selectedImage)
                            .resizable()
                            .aspectRatio(contentMode: .fill)
                            .frame(width: 80, height: 80)
                            .clipShape(Circle())
                    } else {
                        ProfileAvatar(size: 80, user: authService.currentUser)
                    }
                    
                    Circle()
                        .stroke(Color.white, lineWidth: 3)
                        .frame(width: 80, height: 80)
                    
                    VStack {
                        Spacer()
                        HStack {
                            Spacer()
                            Circle()
                                .fill(AppColors.primaryBlue)
                                .frame(width: 24, height: 24)
                                .overlay(
                                    Image(systemName: "camera.fill")
                                        .font(.caption)
                                        .foregroundColor(.white)
                                )
                        }
                    }
                    .frame(width: 80, height: 80)
                }
            }
            
            Text("Tap to change photo")
                .font(AppTextStyles.caption)
                .foregroundColor(.gray)
        }
        .actionSheet(isPresented: $showingActionSheet) {
            ActionSheet(
                title: Text("Select Photo"),
                buttons: [
                    .default(Text("Camera")) {
                        checkCameraPermissionAndShowPicker()
                    },
                    .default(Text("Photo Library")) {
                        checkPhotoLibraryPermissionAndShowPicker()
                    },
                    .cancel()
                ]
            )
        }
        .sheet(isPresented: $showingImagePicker) {
            ImagePicker(selectedImage: $selectedImage, sourceType: imagePickerSourceType)
        }
        .alert("Permission Required", isPresented: $showingPermissionAlert) {
            Button("Settings") {
                if let settingsUrl = URL(string: UIApplication.openSettingsURLString) {
                    UIApplication.shared.open(settingsUrl)
                }
            }
            Button("Cancel", role: .cancel) { }
        } message: {
            Text(permissionMessage)
        }
    }
    
    private func checkCameraPermissionAndShowPicker() {
        PermissionsManager.shared.requestCameraPermission { granted in
            if granted {
                imagePickerSourceType = .camera
                showingImagePicker = true
            } else {
                permissionMessage = "Camera access is required to take photos. Please enable it in Settings."
                showingPermissionAlert = true
            }
        }
    }
    
    private func checkPhotoLibraryPermissionAndShowPicker() {
        PermissionsManager.shared.requestPhotoLibraryPermission { granted in
            if granted {
                imagePickerSourceType = .photoLibrary
                showingImagePicker = true
            } else {
                permissionMessage = "Photo library access is required to select photos. Please enable it in Settings."
                showingPermissionAlert = true
            }
        }
    }
}
