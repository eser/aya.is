import * as runtime from "react/jsx-runtime";
import { compile, runSync } from "@mdx-js/mdx";
import remarkGfm from "remark-gfm";
import rehypeSlug from "rehype-slug";
import rehypeAutolinkHeadings from "rehype-autolink-headings";

// Custom MDX components - demote headings to prevent h1 conflicts
export const mdxComponents = {
  h1: (props: React.HTMLAttributes<HTMLHeadingElement>) => <h2 {...props} />,
  h2: (props: React.HTMLAttributes<HTMLHeadingElement>) => <h3 {...props} />,
  h3: (props: React.HTMLAttributes<HTMLHeadingElement>) => <h4 {...props} />,
  h4: (props: React.HTMLAttributes<HTMLHeadingElement>) => <h5 {...props} />,
  h5: (props: React.HTMLAttributes<HTMLHeadingElement>) => <h6 {...props} />,
};

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
  compiledCode: string
): React.ComponentType<{ components?: Record<string, React.ComponentType> }> {
  const { default: MdxContent } = runSync(compiledCode, {
    ...runtime,
    baseUrl: import.meta.url,
  });

  return MdxContent;
}
