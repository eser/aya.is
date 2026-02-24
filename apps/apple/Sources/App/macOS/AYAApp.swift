import SwiftUI
import Models
import Networking
import UI

@main
struct AYAApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
                .frame(minWidth: 700, minHeight: 500)
        }
        .defaultSize(width: 1000, height: 700)
    }
}
