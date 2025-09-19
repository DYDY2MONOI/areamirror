//
//  IPHelper.swift
//  area
//
//  Created by Dydy2Brazil on 19/09/2025.
//

import Foundation
import Network

class IPHelper: ObservableObject {
    @Published var localIP: String = "localhost"
    
    init() {
        findLocalIP()
    }
    
    private func findLocalIP() {
        let monitor = NWPathMonitor()
        let queue = DispatchQueue(label: "NetworkMonitor")
        
        monitor.pathUpdateHandler = { [weak self] path in
            if path.status == .satisfied {
                self?.getLocalIPAddress()
            }
        }
        
        monitor.start(queue: queue)
    }
    
    private func getLocalIPAddress() {
        var address: String?
        var ifaddr: UnsafeMutablePointer<ifaddrs>?
        
        guard getifaddrs(&ifaddr) == 0 else { return }
        guard let firstAddr = ifaddr else { return }
        
        for ifptr in sequence(first: firstAddr, next: { $0.pointee.ifa_next }) {
            let interface = ifptr.pointee
            
            let addrFamily = interface.ifa_addr.pointee.sa_family
            if addrFamily == UInt8(AF_INET) || addrFamily == UInt8(AF_INET6) {
                
                let name = String(cString: interface.ifa_name)
                if name == "en0" || name == "en1" {
                    
                    var hostname = [CChar](repeating: 0, count: Int(NI_MAXHOST))
                    getnameinfo(interface.ifa_addr, socklen_t(interface.ifa_addr.pointee.sa_len),
                              &hostname, socklen_t(hostname.count),
                              nil, socklen_t(0), NI_NUMERICHOST)
                    address = String(cString: hostname)
                }
            }
        }
        
        freeifaddrs(ifaddr)
        
        DispatchQueue.main.async {
            if let ip = address, !ip.isEmpty {
                self.localIP = ip
            }
        }
    }
    
    func getAPIURL(port: Int = 8080) -> String {
        return "http://\(localIP):\(port)"
    }
}
