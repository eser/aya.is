/**
 * Elements domain - markdown generation utilities
 * Elements shows profiles of individuals and organizations
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/server/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import type { Profile } from "@/modules/backend/types";

/**
 * Generate markdown for the elements listing page
 */
export function generateElementsListingMarkdown(
  profiles: Profile[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Elements",
    generated_at: new Date().toISOString(),
  });

  // Group by kind
  const individuals = profiles.filter((p) => p.kind === "individual");
  const organizations = profiles.filter((p) => p.kind === "organization");

  const sections: string[] = [];

  if (individuals.length > 0) {
    sections.push("## Individuals\n");
    for (const profile of individuals) {
      let line = `- [${profile.title}](/${locale}/${profile.slug}.md)`;
      if (profile.description !== null && profile.description !== undefined) {
        line += `\n  ${profile.description.slice(0, 150)}${profile.description.length > 150 ? "..." : ""}`;
      }
      sections.push(line);
    }
  }

  if (organizations.length > 0) {
    sections.push("\n## Organizations\n");
    for (const profile of organizations) {
      let line = `- [${profile.title}](/${locale}/${profile.slug}.md)`;
      if (profile.description !== null && profile.description !== undefined) {
        line += `\n  ${profile.description.slice(0, 150)}${profile.description.length > 150 ? "..." : ""}`;
      }
      sections.push(line);
    }
  }

  return `${frontmatter}\n\n${sections.join("\n")}`;
}

/**
 * Register markdown handler for elements listing
 * Pattern: /$locale/elements
 */
export function registerElementsListingHandler(): void {
  registerMarkdownHandler("$locale/elements", async (_params, locale, _searchParams) => {
    const profiles = await backend.getProfilesByKinds(locale, [
      "individual",
      "organization",
    ]);

    if (profiles === null) {
      return null;
    }

    return generateElementsListingMarkdown(profiles, locale);
  });
}
