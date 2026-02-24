import SwiftUI
import AYAKit

struct ContentView: View {
    @State private var viewModel = FeedViewModel(client: APIClient())

    var body: some View {
        FeedNavigationView(viewModel: viewModel)
    }
}
