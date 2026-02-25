/**
 * Global middleware to handle markdown requests (.md suffix)
 *
 * Scalable approach:
 * - Each domain registers its markdown handler
 * - Middleware matches URL patterns and delegates to the appropriate handler
 * - Returns markdown Response directly, bypassing React rendering
 *
 * Note: This middleware runs AFTER customDomainMiddleware, so it can use
 * the rewritten path from request context for custom domains.
 */
import { createMiddleware } from "@tanstack/react-start";
import { createMarkdownResponse, createNotFoundResponse } from "@/lib/markdown";
import { DEFAULT_LOCALE } from "@/config";
import { generateLlmsTxt, generateLlmsFullTxt } from "@/routes/-llms-txt";

// Type for markdown handlers
type MarkdownHandler = (
  params: Record<string, string>,
  locale: string,
  searchParams: URLSearchParams,
) => Promise<string | null>;

// Registry of markdown handlers by route pattern
// Pattern format: segments with $ prefix are params (e.g., "$locale/$slug/stories/$storyslug")
const markdownHandlers: Map<string, MarkdownHandler> = new Map();

/**
 * Register a markdown handler for a specific route pattern
 * Pattern uses $ prefix for params: "$locale/$slug/stories/$storyslug"
 */
export function registerMarkdownHandler(
  pattern: string,
  handler: MarkdownHandler,
): void {
  markdownHandlers.set(pattern, handler);
}

/**
 * Calculate pattern specificity (more static segments = more specific)
 */
function getPatternSpecificity(pattern: string): number {
  const segments = pattern.split("/").filter(Boolean);
  return segments.filter((seg) => !seg.startsWith("$")).length;
}

/**
 * Match a URL path against registered patterns
 * Returns all matching handlers sorted by specificity (most specific first)
 */
function matchAllPatterns(
  pathname: string,
): Array<{ handler: MarkdownHandler; params: Record<string, string>; pattern: string }> {
  // Remove .md suffix
  const path = pathname.endsWith(".md") ? pathname.slice(0, -3) : pathname;
  const pathSegments = path.split("/").filter(Boolean);

  const matches: Array<{
    handler: MarkdownHandler;
    params: Record<string, string>;
    pattern: string;
    specificity: number;
  }> = [];

  for (const [pattern, handler] of markdownHandlers) {
    const patternSegments = pattern.split("/").filter(Boolean);

    // Must have same number of segments
    if (pathSegments.length !== patternSegments.length) {
      continue;
    }

    const params: Record<string, string> = {};
    let matched = true;

    for (let i = 0; i < patternSegments.length; i++) {
      const patternSeg = patternSegments[i];
      const pathSeg = pathSegments[i];

      if (patternSeg === undefined || pathSeg === undefined) {
        matched = false;
        break;
      }

      if (patternSeg.startsWith("$")) {
        // This is a param - extract it
        const paramName = patternSeg.slice(1);
        params[paramName] = pathSeg;
      } else if (patternSeg !== pathSeg) {
        // Static segment doesn't match
        matched = false;
        break;
      }
    }

    if (matched) {
      matches.push({
        handler,
        params,
        pattern,
        specificity: getPatternSpecificity(pattern),
      });
    }
  }

  // Sort by specificity (most specific first)
  matches.sort((a, b) => b.specificity - a.specificity);

  return matches;
}

/**
 * Global markdown middleware
 * Intercepts requests ending in .md and returns markdown content
 * Also handles /llms.txt and /llms-full.txt
 *
 * Uses the rewritten path from request context for custom domains.
 */
export const markdownMiddleware = createMiddleware().server(
  async ({ request, next }) => {
    const url = new URL(request.url);

    // Handle llms.txt
    if (url.pathname === "/llms.txt") {
      try {
        const content = await generateLlmsTxt();
        return createMarkdownResponse(content);
      } catch (error) {
        console.error("llms.txt generation error:", error);
        return createNotFoundResponse();
      }
    }

    // Handle llms-full.txt
    if (url.pathname === "/llms-full.txt") {
      try {
        const content = await generateLlmsFullTxt();
        return createMarkdownResponse(content);
      } catch (error) {
        console.error("llms-full.txt generation error:", error);
        return createNotFoundResponse();
      }
    }

    // Only handle .md requests
    if (!url.pathname.endsWith(".md")) {
      return next();
    }

    // Get the effective path - use rewritten path from request context if available
    // This handles custom domains where /tr/index.md becomes /tr/eser/index.md
    const { requestContextBinder } = await import("./request-context-binder");
    const requestContext = requestContextBinder.getStore();
    let effectivePath: string;

    if (requestContext !== undefined && requestContext.path.length > 0) {
      // Use the rewritten path from request context
      // The path array doesn't include the .md suffix, so we need to handle it
      const pathParts = [...requestContext.path];
      const lastPart = pathParts[pathParts.length - 1];
      if (lastPart !== undefined && lastPart.endsWith(".md")) {
        // Last segment has .md suffix
        effectivePath = "/" + pathParts.join("/");
      } else {
        // Append .md suffix from original URL
        effectivePath = "/" + pathParts.join("/") + ".md";
      }
    } else {
      effectivePath = url.pathname;
    }

    // Get all matching handlers sorted by specificity
    const matches = matchAllPatterns(effectivePath);

    if (matches.length === 0) {
      // No handler registered for this pattern
      return createNotFoundResponse();
    }

    // Try each handler in order of specificity until one returns content
    for (const { handler, params } of matches) {
      const locale = params.locale ?? DEFAULT_LOCALE;

      try {
        const markdown = await handler(params, locale, url.searchParams);

        if (markdown !== null) {
          return createMarkdownResponse(markdown, locale);
        }
        // Handler returned null, try next one
      } catch (error) {
        console.error("Markdown generation error:", error);
        // Continue to next handler
      }
    }

    // No handler could produce content
    return createNotFoundResponse();
  },
);
