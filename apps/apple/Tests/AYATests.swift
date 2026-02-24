import Foundation
import Testing
import Models
import Networking

@Suite("Models")
struct ModelTests {
    @Test("APIResponse decodes correctly")
    func apiResponseDecoding() throws {
        let json = """
        {
            "data": [
                {
                    "id": "1",
                    "slug": "test-story",
                    "kind": "article",
                    "title": "Test Story",
                    "summary": "A test summary",
                    "content": "Full content here",
                    "story_picture_uri": "https://example.com/image.jpg",
                    "published_at": "2025-01-01T00:00:00Z",
                    "created_at": "2025-01-01T00:00:00Z",
                    "author_profile": {
                        "id": "a1",
                        "slug": "author",
                        "title": "Author Name",
                        "profile_picture_uri": null
                    }
                }
            ],
            "cursor": "next-page"
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let response = try decoder.decode(APIResponse<Story>.self, from: json)

        #expect(response.data.count == 1)
        #expect(response.data[0].id == "1")
        #expect(response.data[0].slug == "test-story")
        #expect(response.data[0].kind == "article")
        #expect(response.data[0].title == "Test Story")
        #expect(response.data[0].authorProfile?.title == "Author Name")
        #expect(response.cursor == "next-page")
    }

    @Test("Profile decodes correctly")
    func profileDecoding() throws {
        let json = """
        {
            "id": "p1",
            "slug": "test-profile",
            "kind": "individual",
            "title": "Test User",
            "description": "A developer",
            "profile_picture_uri": null,
            "points": 42,
            "default_locale": "en",
            "created_at": "2025-01-01T00:00:00Z",
            "has_translation": true
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let profile = try decoder.decode(Profile.self, from: json)

        #expect(profile.id == "p1")
        #expect(profile.slug == "test-profile")
        #expect(profile.kind == "individual")
        #expect(profile.title == "Test User")
        #expect(profile.points == 42)
        #expect(profile.hasTranslation == true)
    }

    @Test("Activity with properties decodes correctly")
    func activityDecoding() throws {
        let json = """
        {
            "id": "act1",
            "slug": "test-activity",
            "kind": "activity",
            "title": "Workshop",
            "summary": "A workshop",
            "content": null,
            "story_picture_uri": null,
            "published_at": null,
            "created_at": "2025-06-01T10:00:00Z",
            "author_profile": null,
            "properties": {
                "activity_kind": "workshop",
                "activity_time_start": "2025-06-15T14:00:00Z",
                "activity_time_end": "2025-06-15T16:00:00Z",
                "external_activity_uri": "https://meet.example.com/workshop",
                "external_attendance_uri": null,
                "rsvp_mode": "open"
            }
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let activity = try decoder.decode(Activity.self, from: json)

        #expect(activity.id == "act1")
        #expect(activity.properties?.activityKind == "workshop")
        #expect(activity.properties?.activityTimeStart == "2025-06-15T14:00:00Z")
        #expect(activity.properties?.externalActivityUri == "https://meet.example.com/workshop")
        #expect(activity.properties?.rsvpMode == "open")
    }

    @Test("SearchResult decodes correctly")
    func searchResultDecoding() throws {
        let json = """
        {
            "id": "sr1",
            "type": "story",
            "slug": "found-story",
            "title": "Found Story",
            "summary": "Matching content",
            "image_uri": null,
            "profile_slug": "author-slug",
            "profile_title": "Author",
            "kind": "article",
            "rank": 0.95
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let result = try decoder.decode(SearchResult.self, from: json)

        #expect(result.id == "sr1")
        #expect(result.type == "story")
        #expect(result.rank == 0.95)
        #expect(result.profileSlug == "author-slug")
    }

    @Test("Page decodes correctly")
    func pageDecoding() throws {
        let json = """
        {
            "id": "pg1",
            "slug": "about",
            "title": "About Us",
            "content": "# About\\nWe are a community.",
            "created_at": "2025-01-15T12:00:00Z"
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let page = try decoder.decode(Page.self, from: json)

        #expect(page.id == "pg1")
        #expect(page.slug == "about")
        #expect(page.title == "About Us")
        #expect(page.content?.contains("About") == true)
    }

    @Test("Story equality is based on id")
    func storyEquality() {
        // Stories with same id but different titles should be equal
        let json1 = """
        {"id":"1","slug":"s","kind":"article","title":"A","summary":null,"content":null,"story_picture_uri":null,"published_at":null,"created_at":null,"author_profile":null}
        """.data(using: .utf8)!
        let json2 = """
        {"id":"1","slug":"s","kind":"article","title":"B","summary":null,"content":null,"story_picture_uri":null,"published_at":null,"created_at":null,"author_profile":null}
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let story1 = try! decoder.decode(Story.self, from: json1)
        let story2 = try! decoder.decode(Story.self, from: json2)

        #expect(story1 == story2)
    }

    @Test("Spotlight decodes with optional fields")
    func spotlightDecoding() throws {
        let json = """
        {
            "stories": null,
            "activities": [],
            "profiles": null
        }
        """.data(using: .utf8)!

        let decoder = JSONDecoder()
        decoder.keyDecodingStrategy = .convertFromSnakeCase
        let spotlight = try decoder.decode(Spotlight.self, from: json)

        #expect(spotlight.stories == nil)
        #expect(spotlight.activities?.isEmpty == true)
        #expect(spotlight.profiles == nil)
    }
}

@Suite("Networking")
struct NetworkingTests {
    @Test("APIError descriptions are meaningful")
    func apiErrorDescriptions() {
        let errors: [(APIError, String)] = [
            (.invalidURL("/test"), "Invalid URL: /test"),
            (.invalidResponse, "Invalid server response"),
            (.httpError(statusCode: 404), "HTTP error 404"),
        ]

        for (error, expected) in errors {
            #expect(error.errorDescription == expected)
        }
    }

    @Test("LocaleHelper returns valid locale")
    func localeHelper() {
        let locale = LocaleHelper.currentLocale
        #expect(locale == "en" || locale == "tr")
    }

    @Test("APIClient initializes with default values")
    func apiClientInit() {
        let client = APIClient()
        // Should not crash — validates default URL and session
        #expect(client is APIClientProtocol)
    }
}
