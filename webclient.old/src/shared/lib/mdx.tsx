import * as React from "react";

import { compileMDX, type MDXRemoteProps } from "next-mdx-remote/rsc";

export async function mdx(source: string, components?: MDXRemoteProps["components"]) {
  const result = await compileMDX({
    source: source,
    options: {
      mdxOptions: {
        remarkPlugins: [],
        rehypePlugins: [],
        format: "mdx",
      },

      parseFrontmatter: true,
    },
    components: {
      h1: (props) => <h2 {...props} />,
      h2: (props) => <h3 {...props} />,
      h3: (props) => <h4 {...props} />,
      h4: (props) => <h5 {...props} />,
      h5: (props) => <h6 {...props} />,
      ...components,
    },
  });

  return result;
}
