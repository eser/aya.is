import type { ModelContextTool } from "./types.ts";
import { backend } from "@/modules/backend/backend.ts";
import { changeLanguage } from "@/modules/i18n/i18n.ts";
import { isValidLocale, SUPPORTED_LOCALES } from "@/config.ts";

// --- Route context for contextual tools ---

export interface RouteContext {
  profileSlug: string | null;
  storySlug: string | null;
  activitySlug: string | null;
}

const RESERVED_SEGMENTS = new Set([
  "stories",
  "articles",
  "news",
  "activities",
  "contents",
  "products",
  "elements",
  "search",
  "mailbox",
  "admin",
  "auth",
  "api",
  "settings",
]);

export function parseRouteContext(pathname: string): RouteContext {
  const segments = pathname.split("/").filter((s) => s.length > 0);
  const result: RouteContext = {
    profileSlug: null,
    storySlug: null,
    activitySlug: null,
  };

  if (segments.length < 2) {
    return result;
  }

  const second = segments[1];

  // /$locale/stories/$storyslug
  if (second === "stories" && segments.length >= 3) {
    result.storySlug = segments[2];
    return result;
  }

  // /$locale/activities/$activityslug
  if (second === "activities" && segments.length >= 3) {
    result.activitySlug = segments[2];
    return result;
  }

  // /$locale/$slug — profile page (only if not a reserved segment)
  if (!RESERVED_SEGMENTS.has(second)) {
    result.profileSlug = second;

    // /$locale/$slug/stories/$storyslug — story under profile
    if (segments.length >= 4 && segments[2] === "stories") {
      result.storySlug = segments[3];
    }
  }

  return result;
}

// --- Global tools (always available) ---

export function buildGlobalTools(locale: string): ModelContextTool[] {
  return [
    {
      name: "search",
      description:
        "Search for profiles (people, organizations, products) and content on AYA. Returns matching results with titles, summaries, and slugs.",
      inputSchema: {
        type: "object",
        properties: {
          query: { type: "string", description: "Search query text" },
          limit: {
            type: "number",
            description: "Maximum results to return (default: 20)",
          },
        },
        required: ["query"],
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const query = input.query as string;
        const limit = typeof input.limit === "number" ? input.limit : 20;
        return await backend.search(locale, query, undefined, limit);
      },
    },

    {
      name: "get-profile",
      description:
        "Get detailed information about a profile (person, organization, or product) including description, links, and pages.",
      inputSchema: {
        type: "object",
        properties: {
          slug: {
            type: "string",
            description: "Profile URL slug (e.g., 'eser')",
          },
        },
        required: ["slug"],
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const slug = input.slug as string;
        return await backend.getProfile(locale, slug);
      },
    },

    {
      name: "get-profile-page",
      description: "Get the content of a custom page belonging to a profile (e.g., an 'about' page or CV).",
      inputSchema: {
        type: "object",
        properties: {
          profile_slug: {
            type: "string",
            description: "Profile URL slug",
          },
          page_slug: {
            type: "string",
            description: "Page URL slug within the profile",
          },
        },
        required: ["profile_slug", "page_slug"],
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const profileSlug = input.profile_slug as string;
        const pageSlug = input.page_slug as string;
        return await backend.getProfilePage(locale, profileSlug, pageSlug);
      },
    },

    {
      name: "get-story",
      description:
        "Get the full content of a story, article, or news item. Returns title, summary, content, author info, and metadata.",
      inputSchema: {
        type: "object",
        properties: {
          slug: { type: "string", description: "Story URL slug" },
        },
        required: ["slug"],
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const slug = input.slug as string;
        return await backend.getStory(locale, slug);
      },
    },

    {
      name: "get-story-cover-image",
      description: "Get the cover image URL for a story or article.",
      inputSchema: {
        type: "object",
        properties: {
          slug: { type: "string", description: "Story URL slug" },
        },
        required: ["slug"],
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const slug = input.slug as string;
        const story = await backend.getStory(locale, slug);
        if (story === null) {
          return null;
        }
        return {
          url: story.story_picture_uri,
          title: story.title,
          slug: story.slug,
        };
      },
    },

    {
      name: "list-stories",
      description:
        "List recent stories, optionally filtered by kind. Kinds: article, announcement, news, status, content, presentation, activity.",
      inputSchema: {
        type: "object",
        properties: {
          kind: {
            type: "string",
            description: "Filter by story kind (e.g., 'article', 'news'). Omit for all kinds.",
            enum: [
              "article",
              "announcement",
              "news",
              "status",
              "content",
              "presentation",
              "activity",
            ],
          },
        },
      },
      annotations: { readOnlyHint: true },
      execute: async (input: Record<string, unknown>) => {
        const kind = input.kind as string | undefined;
        if (kind !== undefined && kind !== "") {
          return await backend.getStoriesByKinds(locale, [kind]);
        }
        return await backend.getStories(locale);
      },
    },

    {
      name: "list-activities",
      description: "List upcoming activities and events including workshops, conferences, webinars, and broadcasts.",
      inputSchema: { type: "object", properties: {} },
      annotations: { readOnlyHint: true },
      execute: async () => {
        return await backend.getActivities(locale);
      },
    },

    {
      name: "navigate-to",
      description:
        "Navigate the browser to a specific page on AYA. Provide a path like '/{locale}/{slug}' or '/{locale}/stories/{slug}'.",
      inputSchema: {
        type: "object",
        properties: {
          path: {
            type: "string",
            description: "URL path to navigate to (must start with '/')",
          },
        },
        required: ["path"],
      },
      annotations: { readOnlyHint: false },
      execute: (input: Record<string, unknown>) => {
        const path = input.path as string;
        if (!path.startsWith("/")) {
          return Promise.reject(
            new Error(
              "Path must start with '/' to prevent external navigation.",
            ),
          );
        }
        globalThis.location.href = path;
        return Promise.resolve({ navigated: true, path });
      },
    },

    {
      name: "switch-language",
      description: `Switch the site's display language. Supported locales: ${SUPPORTED_LOCALES.join(", ")}.`,
      inputSchema: {
        type: "object",
        properties: {
          locale: {
            type: "string",
            description: "Target locale code",
            enum: [...SUPPORTED_LOCALES],
          },
        },
        required: ["locale"],
      },
      annotations: { readOnlyHint: false },
      execute: async (input: Record<string, unknown>) => {
        const targetLocale = input.locale as string;
        if (!isValidLocale(targetLocale)) {
          throw new Error(
            `Unsupported locale: ${targetLocale}. Supported: ${SUPPORTED_LOCALES.join(", ")}`,
          );
        }
        await changeLanguage(targetLocale);
        return { switched: true, locale: targetLocale };
      },
    },

    {
      name: "copy-story-markdown",
      description: "Copy a story's content as markdown to the user's clipboard.",
      inputSchema: {
        type: "object",
        properties: {
          slug: { type: "string", description: "Story URL slug" },
        },
        required: ["slug"],
      },
      annotations: { readOnlyHint: false },
      execute: async (input: Record<string, unknown>) => {
        const slug = input.slug as string;
        const response = await fetch(`/${locale}/stories/${slug}.md`);
        if (!response.ok) {
          throw new Error(`Story not found: ${slug}`);
        }
        const markdown = await response.text();
        await navigator.clipboard.writeText(markdown);
        return { copied: true, slug, length: markdown.length };
      },
    },
  ];
}

// --- Contextual tools (page-aware, added based on current route) ---

export function buildContextualTools(
  locale: string,
  context: RouteContext,
): ModelContextTool[] {
  const tools: ModelContextTool[] = [];

  if (context.profileSlug !== null) {
    const slug = context.profileSlug;
    tools.push({
      name: "get-current-profile",
      description: `Get details about the profile currently being viewed (${slug}).`,
      inputSchema: { type: "object", properties: {} },
      annotations: { readOnlyHint: true },
      execute: async () => {
        return await backend.getProfile(locale, slug);
      },
    });
  }

  if (context.storySlug !== null) {
    const slug = context.storySlug;
    tools.push(
      {
        name: "get-current-story",
        description: `Get the content of the story currently being viewed (${slug}).`,
        inputSchema: { type: "object", properties: {} },
        annotations: { readOnlyHint: true },
        execute: async () => {
          return await backend.getStory(locale, slug);
        },
      },
      {
        name: "copy-current-story-markdown",
        description: `Copy the current story (${slug}) as markdown to the clipboard.`,
        inputSchema: { type: "object", properties: {} },
        annotations: { readOnlyHint: false },
        execute: async () => {
          const response = await fetch(
            `/${locale}/stories/${slug}.md`,
          );
          if (!response.ok) {
            throw new Error(`Story not found: ${slug}`);
          }
          const markdown = await response.text();
          await navigator.clipboard.writeText(markdown);
          return { copied: true, slug, length: markdown.length };
        },
      },
    );
  }

  if (context.activitySlug !== null) {
    const slug = context.activitySlug;
    tools.push({
      name: "get-current-activity",
      description: `Get details about the activity currently being viewed (${slug}).`,
      inputSchema: { type: "object", properties: {} },
      annotations: { readOnlyHint: true },
      execute: async () => {
        return await backend.getStory(locale, slug);
      },
    });
  }

  return tools;
}
