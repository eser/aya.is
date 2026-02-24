// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "AYA",
    platforms: [
        .iOS(.v18),
        .macOS(.v15),
    ],
    products: [
        .library(name: "AYAKit", targets: ["AYAKit"]),
    ],
    targets: [
        .target(
            name: "AYAKit",
            path: "Sources",
            exclude: ["App"]
        ),
        .testTarget(
            name: "AYATests",
            dependencies: [
                "AYAKit",
            ],
            path: "Tests"
        ),
    ]
)
