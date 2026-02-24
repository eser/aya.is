.PHONY: project build-ios build-macos test test-ios test-macos clean lint \
       fastlane-ios-build fastlane-macos-build fastlane-ios-test fastlane-macos-test \
       fastlane-setup check format

# --- Project Generation ---

project:
	xcodegen generate

# --- Build ---

build-ios: project
	xcodebuild build \
		-project AYA.xcodeproj \
		-scheme AYA-iOS \
		-destination 'platform=iOS Simulator,name=iPhone 16' \
		-skipPackagePluginValidation \
		CODE_SIGNING_ALLOWED=NO

build-macos: project
	xcodebuild build \
		-project AYA.xcodeproj \
		-scheme AYA-macOS \
		-destination 'platform=macOS' \
		-skipPackagePluginValidation \
		CODE_SIGNING_ALLOWED=NO

build: build-ios build-macos

# --- Tests ---

test:
	swift test

test-ios: project
	xcodebuild test \
		-project AYA.xcodeproj \
		-scheme AYA-iOS \
		-destination 'platform=iOS Simulator,name=iPhone 16' \
		-skipPackagePluginValidation \
		CODE_SIGNING_ALLOWED=NO

test-macos: project
	xcodebuild test \
		-project AYA.xcodeproj \
		-scheme AYA-macOS \
		-destination 'platform=macOS' \
		-skipPackagePluginValidation \
		CODE_SIGNING_ALLOWED=NO

# --- Lint ---

lint:
	@echo "Checking for removed DI patterns..."
	@! grep -r "ComponentResolver\|ComponentBinding\|UIFactory\|service.locator\|AppLogger" Sources/ || \
		(echo "ERROR: Found references to removed DI patterns" && exit 1)
	@echo "Clean — no legacy patterns found."

format:
	swift format --in-place --recursive Sources/ Tests/

check: lint test
	@echo "All checks passed."

# --- Clean ---

clean:
	xcodebuild clean -project AYA.xcodeproj -scheme AYA-iOS 2>/dev/null || true
	xcodebuild clean -project AYA.xcodeproj -scheme AYA-macOS 2>/dev/null || true
	rm -rf .build
	rm -rf ~/Library/Developer/Xcode/DerivedData/AYA-*

# --- Fastlane ---

fastlane-setup:
	bundle install

fastlane-ios-build:
	bundle exec fastlane ios build

fastlane-macos-build:
	bundle exec fastlane mac build

fastlane-ios-test:
	bundle exec fastlane ios test

fastlane-macos-test:
	bundle exec fastlane mac test

fastlane-ios-archive:
	bundle exec fastlane ios archive

fastlane-macos-archive:
	bundle exec fastlane mac archive

fastlane-beta:
	bundle exec fastlane ios beta
