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
    dependencies: [
        .package(url: "https://github.com/pointfreeco/swift-snapshot-testing", from: "1.17.0"),
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
                .product(name: "SnapshotTesting", package: "swift-snapshot-testing"),
            ],
            path: "Tests"
        ),
    ]
)
