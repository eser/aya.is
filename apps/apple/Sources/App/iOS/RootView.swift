import SwiftUI
import AYAKit

struct RootView: View {
    @State private var viewModel: FeedViewModel

    init() {
        if ProcessInfo.processInfo.arguments.contains("--uitesting") {
            let session = MockURLProtocol.mockSession()
            _viewModel = State(initialValue: FeedViewModel(client: APIClient(session: session)))
        } else {
            _viewModel = State(initialValue: FeedViewModel(client: APIClient()))
        }
    }

    var body: some View {
        FeedNavigationView(viewModel: viewModel)
    }
}
