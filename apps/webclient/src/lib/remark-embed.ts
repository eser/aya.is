import { visit } from "unist-util-visit";
import type { Root, Paragraph, Text, Link } from "mdast";

/**
 * A remark plugin that transforms %[url] syntax into <Embed url="..." /> components.
 * This is compatible with Hashnode's embed syntax.
 *
 * Example:
 *   %[https://www.youtube.com/watch?v=abc123]
 *   becomes:
 *   <Embed url="https://www.youtube.com/watch?v=abc123" />
 */
export function remarkEmbed() {
  return (tree: Root) => {
    visit(tree, "paragraph", (node: Paragraph, index, parent) => {
      if (parent === undefined || parent === null || typeof index !== "number") {
        return;
      }

      let url: string | null = null;

      // Case 1: Single text node with %[url] pattern (before autolink conversion)
      if (node.children.length === 1 && node.children[0].type === "text") {
        const text = (node.children[0] as Text).value.trim();
        const match = text.match(/^%\[([^\]]+)\]$/);
        if (match !== null) {
          url = match[1];
        }
      }

      // Case 2: Pattern was split by autolink: text "%[" + link + text "]"
      if (url === null && node.children.length === 3) {
        const [first, second, third] = node.children;

        if (
          first.type === "text" &&
          second.type === "link" &&
          third.type === "text"
        ) {
          const prefix = (first as Text).value;
          const suffix = (third as Text).value;

          if (prefix === "%[" && suffix === "]") {
            url = (second as Link).url;
          }
        }
      }

      // Case 3: Just a link inside %[] where text nodes might be empty or combined
      if (url === null && node.children.length >= 1) {
        // Build full text from all children
        let fullText = "";
        let foundLink: Link | null = null;

        for (const child of node.children) {
          if (child.type === "text") {
            fullText += (child as Text).value;
          } else if (child.type === "link" && foundLink === null) {
            foundLink = child as Link;
            fullText += `__LINK__`;
          }
        }

        const trimmed = fullText.trim();
        if (trimmed === "%[__LINK__]" && foundLink !== null) {
          url = foundLink.url;
        }
      }

      if (url === null) {
        return;
      }

      // Replace the paragraph with an MDX JSX element
      parent.children[index] = {
        type: "mdxJsxFlowElement",
        name: "Embed",
        attributes: [
          {
            type: "mdxJsxAttribute",
            name: "url",
            value: url,
          },
        ],
        children: [],
        data: { _mdxExplicitJsx: true },
      } as any;
    });
  };
}
