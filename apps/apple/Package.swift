// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "AYA",
    platforms: [
        .iOS(.v18),
        .macOS(.v15),
    ],
    products: [
        .library(name: "AYA", targets: ["Models", "Networking", "UI"]),
    ],
    targets: [
        .target(
            name: "Models",
            path: "Sources/Models"
        ),
        .target(
            name: "Networking",
            dependencies: ["Models"],
            path: "Sources/Networking"
        ),
        .target(
            name: "UI",
            dependencies: ["Models", "Networking"],
            path: "Sources/UI"
        ),
        .testTarget(
            name: "AYATests",
            dependencies: ["Models", "Networking", "UI"],
            path: "Tests"
        ),
    ]
)
