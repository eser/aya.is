import * as runtime from "react/jsx-runtime";
import { compile, runSync } from "@mdx-js/mdx";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeAutolinkHeadings from "rehype-autolink-headings";
import { visit } from "unist-util-visit";
import { remarkEmbed } from "./remark-embed";

import {
  Card,
  Cards,
  Embed,
  List,
  ListItem,
  ProfileCard,
  ProfileList,
  SiteLink,
  Story,
  TwitterEmbed,
  YouTubeEmbed,
} from "@/components/userland";

type HeadingTag = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

const headingTags: HeadingTag[] = ["h1", "h2", "h3", "h4", "h5", "h6"];
/**
 * Rehype plugin that adds target="_blank" and rel="noopener noreferrer"
 * to all links so they open in a new tab.
 */
function rehypeExternalLinks() {
  return (tree: any) => {
    visit(tree, "element", (node: any) => {
      if (
        node.tagName === "a" &&
        node.properties?.href !== undefined &&
        typeof node.properties.href === "string" &&
        URL.canParse(node.properties.href)
      ) {
        node.properties.target = "_blank";
        node.properties.rel = "noopener noreferrer";
      }
    });
  };
}

/**
 * Creates MDX components with heading level offset.
 * offset=1 means h1→h2, h2→h3, etc. (default, for h1 page titles)
 * offset=2 means h1→h3, h2→h4, etc. (for h2 page titles)
 *
 * When includeUserlandComponents is false, custom components (Card, Embed, etc.)
 * are excluded. Useful for user-generated content like Q&A where custom
 * MDX components should not be rendered.
 */
export function createMdxComponents(
  offset: number = 1,
  includeUserlandComponents: boolean = true,
): Record<string, React.ComponentType<any>> {
  const components: Record<string, React.ComponentType<any>> = {};

  // Add heading components with offset
  for (let i = 0; i < headingTags.length; i++) {
    const sourceTag = headingTags[i];
    const targetIndex = Math.min(i + offset, headingTags.length - 1);
    const TargetTag = headingTags[targetIndex];

    components[sourceTag] = (props: React.HTMLAttributes<HTMLHeadingElement>) => (
      <TargetTag {...props} />
    );
  }

  if (includeUserlandComponents) {
    // Add custom userland components
    components.Card = Card;
    components.Cards = Cards;
    components.Embed = Embed;
    components.List = List;
    components.ListItem = ListItem;
    components.ProfileCard = ProfileCard;
    components.ProfileList = ProfileList;
    components.SiteLink = SiteLink;
    components.Story = Story;
    components.TwitterEmbed = TwitterEmbed;
    components.YouTubeEmbed = YouTubeEmbed;
  }

  return components;
}

// Default components with offset 1 (for backward compatibility)
export const mdxComponents = createMdxComponents(1);

/**
 * Compiles MDX source to a function body string that can be serialized.
 * This runs on the server (in the loader).
 */
export async function compileMdx(source: string): Promise<string> {
  const compiled = await compile(source, {
    outputFormat: "function-body",
    remarkPlugins: [remarkEmbed, remarkGfm],
    rehypePlugins: [
      rehypeSlug,
      [rehypeAutolinkHeadings, { behavior: "wrap" }],
      rehypeExternalLinks,
    ],
  });

  return String(compiled);
}

/**
 * Runs compiled MDX code synchronously and returns the React component.
 * This can run on both server and client for SSR support.
 */
export function runMdxSync(
  compiledCode: string,
): React.ComponentType<{ components?: Record<string, React.ComponentType> }> {
  const { default: MdxContent } = runSync(compiledCode, {
    ...runtime,
    baseUrl: import.meta.url,
  });

  return MdxContent;
}
