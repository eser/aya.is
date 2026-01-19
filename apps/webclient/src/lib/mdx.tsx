import * as runtime from "react/jsx-runtime";
import { compile, runSync } from "@mdx-js/mdx";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeAutolinkHeadings from "rehype-autolink-headings";
import {
  Card,
  Cards,
  List,
  ListItem,
  ProfileCard,
  ProfileList,
  SiteLink,
  Story,
} from "@/components/userland";

type HeadingTag = "h1" | "h2" | "h3" | "h4" | "h5" | "h6";

const headingTags: HeadingTag[] = ["h1", "h2", "h3", "h4", "h5", "h6"];

/**
 * Creates MDX components with heading level offset.
 * offset=1 means h1→h2, h2→h3, etc. (default, for h1 page titles)
 * offset=2 means h1→h3, h2→h4, etc. (for h2 page titles)
 */
export function createMdxComponents(
  offset: number = 1,
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

  // Add custom userland components
  components.Card = Card;
  components.Cards = Cards;
  components.List = List;
  components.ListItem = ListItem;
  components.ProfileCard = ProfileCard;
  components.ProfileList = ProfileList;
  components.SiteLink = SiteLink;
  components.Story = Story;

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
    remarkPlugins: [remarkGfm],
    rehypePlugins: [
      rehypeSlug,
      [rehypeAutolinkHeadings, { behavior: "wrap" }],
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
