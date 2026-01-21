/**
 * Central registration of all markdown handlers
 *
 * Each domain exports a registration function that is called here.
 * This file is imported by start.ts to ensure handlers are registered
 * before the middleware runs.
 */

// Domain-specific markdown handler registrations
// Files prefixed with - are excluded from TanStack Router's route scanning
import { registerProfileStoryMarkdownHandler } from "@/routes/$locale/$slug/stories/-markdown";
import {
  registerProfileMarkdownHandler,
  registerProfilePageMarkdownHandler,
} from "@/routes/$locale/$slug/-markdown";
import {
  registerGlobalStoriesListingHandler,
  registerGlobalStoryMarkdownHandler,
} from "@/routes/$locale/stories/-markdown";
import { registerNewsListingHandler } from "@/routes/$locale/news/-markdown";
import { registerProductsListingHandler } from "@/routes/$locale/products/-markdown";
import { registerElementsListingHandler } from "@/routes/$locale/elements/-markdown";
import { registerEventsListingHandler } from "@/routes/$locale/events/-markdown";
import { registerLocaleIndexHandler } from "@/routes/$locale/-markdown";

/**
 * Register all markdown handlers
 * Called once during app initialization
 */
export function registerAllMarkdownHandlers(): void {
  // Profile stories: /$locale/$slug/stories/$storyslug
  registerProfileStoryMarkdownHandler();

  // Profile pages: /$locale/$slug
  registerProfileMarkdownHandler();

  // Profile custom pages: /$locale/$slug/$pageslug
  registerProfilePageMarkdownHandler();

  // Global stories listing: /$locale/stories
  registerGlobalStoriesListingHandler();

  // Global stories: /$locale/stories/$storyslug
  registerGlobalStoryMarkdownHandler();

  // News listing: /$locale/news
  registerNewsListingHandler();

  // Products listing: /$locale/products
  registerProductsListingHandler();

  // Elements listing: /$locale/elements
  registerElementsListingHandler();

  // Events listing: /$locale/events
  registerEventsListingHandler();

  // Locale index: /$locale
  registerLocaleIndexHandler();
}
