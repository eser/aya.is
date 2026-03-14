import type { BlockCategory, BlockDefinition, BlockPattern } from "./types.ts";
import { generateContainerMdx, generateSelfClosingMdx } from "./mdx-template.ts";
import {
  AlertTriangle,
  ChevronDown,
  Code,
  Columns2,
  Columns3,
  CreditCard,
  ExternalLink,
  FileText,
  GitCompareArrows,
  Grid3X3,
  Heading2,
  HelpCircle,
  Image,
  LayoutGrid,
  LayoutTemplate,
  Link,
  Megaphone,
  MessageSquareQuote,
  Minus,
  MousePointerClick,
  PanelTop,
  Quote,
  Space,
  Twitter,
  Youtube,
} from "lucide-react";

// ---------------------------------------------------------------------------
// Block definitions
// ---------------------------------------------------------------------------

const blockDefinitions: BlockDefinition[] = [
  // ---- text ----
  {
    id: "callout",
    name: "Callout",
    description: "Highlight important information with a styled callout box",
    icon: AlertTriangle,
    category: "text",
    keywords: ["callout", "alert", "info", "warning", "tip", "danger", "note"],
    props: [
      {
        name: "variant",
        type: "select",
        label: "Variant",
        required: false,
        defaultValue: "info",
        options: [
          { value: "info", label: "Info" },
          { value: "warning", label: "Warning" },
          { value: "tip", label: "Tip" },
          { value: "danger", label: "Danger" },
        ],
      },
      {
        name: "title",
        type: "string",
        label: "Title",
        required: false,
        defaultValue: "",
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.variant !== undefined && values.variant !== "info") {
        props.variant = values.variant;
      }
      if (typeof values.title === "string" && values.title !== "") {
        props.title = values.title;
      }
      return generateContainerMdx("Callout", props, "Your content here.");
    },
  },
  {
    id: "details",
    name: "Details",
    description: "Collapsible section with a summary toggle",
    icon: ChevronDown,
    category: "text",
    keywords: ["details", "collapse", "accordion", "expand", "toggle"],
    props: [
      {
        name: "summary",
        type: "string",
        label: "Summary",
        required: true,
        defaultValue: "Click to expand",
      },
    ],
    generateMdx: (values) => {
      const summary = typeof values.summary === "string" && values.summary !== "" ? values.summary : "Click to expand";
      return generateContainerMdx(
        "Details",
        { summary },
        "Hidden content here.",
      );
    },
  },
  {
    id: "heading",
    name: "Heading",
    description: "Section heading (h2, h3, or h4)",
    icon: Heading2,
    category: "text",
    keywords: ["heading", "title", "h2", "h3", "h4", "section"],
    props: [
      {
        name: "level",
        type: "select",
        label: "Level",
        required: false,
        defaultValue: "2",
        options: [
          { value: "2", label: "H2" },
          { value: "3", label: "H3" },
          { value: "4", label: "H4" },
        ],
      },
      {
        name: "text",
        type: "string",
        label: "Text",
        required: false,
        defaultValue: "Heading text",
      },
    ],
    generateMdx: (values) => {
      const level = typeof values.level === "string" ? Number(values.level) : 2;
      const text = typeof values.text === "string" && values.text !== "" ? values.text : "Heading text";
      const hashes = "#".repeat(level);
      return `\n${hashes} ${text}\n`;
    },
  },
  {
    id: "quote",
    name: "Quote",
    description: "Block quote for citations or callouts",
    icon: Quote,
    category: "text",
    keywords: ["quote", "blockquote", "citation", "cite"],
    props: [
      {
        name: "text",
        type: "string",
        label: "Text",
        required: false,
        defaultValue: "Quote text",
      },
    ],
    generateMdx: (values) => {
      const text = typeof values.text === "string" && values.text !== "" ? values.text : "Quote text";
      return `\n> ${text}\n`;
    },
  },
  {
    id: "code-block",
    name: "Code Block",
    description: "Syntax-highlighted code snippet",
    icon: Code,
    category: "text",
    keywords: ["code", "snippet", "syntax", "highlight", "pre"],
    props: [
      {
        name: "language",
        type: "string",
        label: "Language",
        required: false,
        defaultValue: "js",
      },
      {
        name: "code",
        type: "string",
        label: "Code",
        required: false,
        defaultValue: "code here",
      },
    ],
    generateMdx: (values) => {
      const language = typeof values.language === "string" && values.language !== "" ? values.language : "js";
      const code = typeof values.code === "string" && values.code !== "" ? values.code : "code here";
      return `\n\`\`\`${language}\n${code}\n\`\`\`\n`;
    },
  },

  // ---- media ----
  {
    id: "image",
    name: "Image",
    description: "Display an image with alt text",
    icon: Image,
    category: "media",
    keywords: ["image", "photo", "picture", "img"],
    props: [
      {
        name: "src",
        type: "string",
        label: "Source URL",
        required: true,
        defaultValue: "",
      },
      {
        name: "alt",
        type: "string",
        label: "Alt Text",
        required: false,
        defaultValue: "Alt text",
      },
    ],
    generateMdx: (values) => {
      const src = typeof values.src === "string" ? values.src : "";
      const alt = typeof values.alt === "string" && values.alt !== "" ? values.alt : "Alt text";
      return `\n![${alt}](${src})\n`;
    },
  },

  // ---- layout ----
  {
    id: "columns-2",
    name: "2 Columns",
    description: "Two-column layout",
    icon: Columns2,
    category: "layout",
    keywords: ["columns", "two", "2", "grid", "side by side", "layout"],
    props: [],
    generateMdx: () => {
      return generateContainerMdx(
        "Columns",
        { cols: 2 },
        `<Column>\n\nColumn 1 content\n\n</Column>\n\n<Column>\n\nColumn 2 content\n\n</Column>`,
      );
    },
  },
  {
    id: "columns-3",
    name: "3 Columns",
    description: "Three-column layout",
    icon: Columns3,
    category: "layout",
    keywords: ["columns", "three", "3", "grid", "layout"],
    props: [],
    generateMdx: () => {
      return generateContainerMdx(
        "Columns",
        { cols: 3 },
        `<Column>\n\nColumn 1 content\n\n</Column>\n\n<Column>\n\nColumn 2 content\n\n</Column>\n\n<Column>\n\nColumn 3 content\n\n</Column>`,
      );
    },
  },
  {
    id: "divider",
    name: "Divider",
    description: "Visual separator between content sections",
    icon: Minus,
    category: "layout",
    keywords: ["divider", "separator", "line", "hr", "break"],
    props: [
      {
        name: "variant",
        type: "select",
        label: "Variant",
        required: false,
        defaultValue: "line",
        options: [
          { value: "line", label: "Line" },
          { value: "dots", label: "Dots" },
          { value: "space", label: "Space" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.variant !== undefined && values.variant !== "line") {
        props.variant = values.variant;
      }
      return generateSelfClosingMdx("Divider", props);
    },
  },
  {
    id: "spacer",
    name: "Spacer",
    description: "Vertical spacing between content sections",
    icon: Space,
    category: "layout",
    keywords: ["spacer", "space", "gap", "padding", "margin"],
    props: [
      {
        name: "size",
        type: "select",
        label: "Size",
        required: false,
        defaultValue: "md",
        options: [
          { value: "sm", label: "Small" },
          { value: "md", label: "Medium" },
          { value: "lg", label: "Large" },
          { value: "xl", label: "Extra Large" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.size !== undefined && values.size !== "md") {
        props.size = values.size;
      }
      return generateSelfClosingMdx("Spacer", props);
    },
  },
  {
    id: "cards",
    name: "Cards",
    description: "Grid container for Card components",
    icon: LayoutGrid,
    category: "layout",
    keywords: ["cards", "grid", "card grid", "card list"],
    props: [],
    generateMdx: () => {
      return `\n<Cards>\n\n<Card title="Card 1" description="Description" />\n\n<Card title="Card 2" description="Description" />\n\n</Cards>\n`;
    },
  },

  // ---- data ----
  {
    id: "card",
    name: "Card",
    description: "Content card with title, description, and optional link",
    icon: CreditCard,
    category: "data",
    keywords: ["card", "box", "tile", "panel"],
    props: [
      {
        name: "title",
        type: "string",
        label: "Title",
        required: true,
        defaultValue: "Title",
      },
      {
        name: "description",
        type: "string",
        label: "Description",
        required: true,
        defaultValue: "Description",
      },
      {
        name: "href",
        type: "string",
        label: "Link URL",
        required: false,
        defaultValue: "#",
      },
    ],
    generateMdx: (values) => {
      const title = typeof values.title === "string" && values.title !== "" ? values.title : "Title";
      const description = typeof values.description === "string" && values.description !== ""
        ? values.description
        : "Description";
      const props: Record<string, string | number | boolean> = {
        title,
        description,
      };
      if (typeof values.href === "string" && values.href !== "" && values.href !== "#") {
        props.href = values.href;
      }
      return generateSelfClosingMdx("Card", props);
    },
  },
  {
    id: "site-link",
    name: "Site Link",
    description: "Internal or external link with custom text",
    icon: ExternalLink,
    category: "data",
    keywords: ["link", "site link", "internal link", "navigation"],
    props: [
      {
        name: "href",
        type: "string",
        label: "URL",
        required: true,
        defaultValue: "/page",
      },
      {
        name: "text",
        type: "string",
        label: "Link Text",
        required: false,
        defaultValue: "Link text",
      },
    ],
    generateMdx: (values) => {
      const href = typeof values.href === "string" && values.href !== "" ? values.href : "/page";
      const text = typeof values.text === "string" && values.text !== "" ? values.text : "Link text";
      return generateContainerMdx("SiteLink", { href }, text);
    },
  },

  // ---- embed ----
  {
    id: "youtube",
    name: "YouTube",
    description: "Embed a YouTube video",
    icon: Youtube,
    category: "embed",
    keywords: ["youtube", "video", "embed", "player"],
    props: [
      {
        name: "videoId",
        type: "string",
        label: "Video ID",
        required: true,
        defaultValue: "",
        placeholder: "dQw4w9WgXcQ",
      },
    ],
    generateMdx: (values) => {
      const videoId = typeof values.videoId === "string" ? values.videoId : "";
      return generateSelfClosingMdx("YouTube", { videoId });
    },
  },
  {
    id: "twitter",
    name: "Twitter",
    description: "Embed a tweet",
    icon: Twitter,
    category: "embed",
    keywords: ["twitter", "tweet", "x", "social", "embed"],
    props: [
      {
        name: "tweetId",
        type: "string",
        label: "Tweet ID",
        required: true,
        defaultValue: "",
      },
      {
        name: "username",
        type: "string",
        label: "Username",
        required: true,
        defaultValue: "",
      },
    ],
    generateMdx: (values) => {
      const tweetId = typeof values.tweetId === "string" ? values.tweetId : "";
      const username = typeof values.username === "string" ? values.username : "";
      return generateSelfClosingMdx("Tweet", { tweetId, username });
    },
  },
  {
    id: "pdf",
    name: "PDF",
    description: "Embed a PDF document viewer",
    icon: FileText,
    category: "embed",
    keywords: ["pdf", "document", "viewer", "embed"],
    props: [
      {
        name: "src",
        type: "string",
        label: "PDF URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com/doc.pdf",
      },
    ],
    generateMdx: (values) => {
      const src = typeof values.src === "string" ? values.src : "";
      return generateSelfClosingMdx("PDF", { src });
    },
  },
  {
    id: "embed",
    name: "Embed",
    description: "Embed external content via URL",
    icon: Link,
    category: "embed",
    keywords: ["embed", "iframe", "external", "url"],
    props: [
      {
        name: "url",
        type: "string",
        label: "URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com",
      },
    ],
    generateMdx: (values) => {
      const url = typeof values.url === "string" ? values.url : "";
      return generateSelfClosingMdx("Embed", { url });
    },
  },

  // ---- interactive ----
  {
    id: "button",
    name: "Button",
    description: "Call-to-action button with a link",
    icon: MousePointerClick,
    category: "interactive",
    keywords: ["button", "cta", "call to action", "link", "action"],
    props: [
      {
        name: "href",
        type: "string",
        label: "URL",
        required: true,
        defaultValue: "",
      },
      {
        name: "variant",
        type: "select",
        label: "Variant",
        required: false,
        defaultValue: "default",
        options: [
          { value: "default", label: "Default" },
          { value: "outline", label: "Outline" },
          { value: "secondary", label: "Secondary" },
        ],
      },
    ],
    generateMdx: (values) => {
      const href = typeof values.href === "string" ? values.href : "";
      const props: Record<string, string | number | boolean> = { href };
      if (values.variant !== undefined && values.variant !== "default") {
        props.variant = values.variant;
      }
      return generateContainerMdx("Button", props, "Click here");
    },
  },
  {
    id: "tabs",
    name: "Tabs",
    description: "Tabbed content sections",
    icon: PanelTop,
    category: "interactive",
    keywords: ["tabs", "tab", "tabbed", "switch", "panel"],
    props: [],
    generateMdx: () => {
      return generateContainerMdx(
        "Tabs",
        {},
        `<Tab label="Tab 1">\n\nTab 1 content\n\n</Tab>\n\n<Tab label="Tab 2">\n\nTab 2 content\n\n</Tab>`,
      );
    },
  },
];

// ---------------------------------------------------------------------------
// Block patterns
// ---------------------------------------------------------------------------

const blockPatterns: BlockPattern[] = [
  {
    id: "hero-section",
    name: "Hero Section",
    description: "A hero area with heading, image, and call-to-action button",
    icon: LayoutTemplate,
    category: "section",
    template: `## Hero Title

![Hero image](https://placehold.co/1200x400)

<Button href="/get-started">

Get Started

</Button>
`,
  },
  {
    id: "feature-grid",
    name: "Feature Grid",
    description: "A grid of feature cards highlighting key points",
    icon: Grid3X3,
    category: "section",
    template: `<Cards>

<Card title="Feature 1" description="Description of the first feature" />

<Card title="Feature 2" description="Description of the second feature" />

<Card title="Feature 3" description="Description of the third feature" />

</Cards>
`,
  },
  {
    id: "faq",
    name: "FAQ",
    description: "Frequently asked questions with collapsible answers",
    icon: HelpCircle,
    category: "section",
    template: `<Details summary="What is this?">

Answer to the first question goes here.

</Details>

<Details summary="How does it work?">

Answer to the second question goes here.

</Details>

<Details summary="Where can I learn more?">

Answer to the third question goes here.

</Details>
`,
  },
  {
    id: "testimonial",
    name: "Testimonial",
    description: "A customer testimonial with quote and attribution link",
    icon: MessageSquareQuote,
    category: "section",
    template: `> This product completely transformed how our team collaborates. Highly recommended!

<SiteLink href="/about">

-- Jane Doe, Acme Corp

</SiteLink>
`,
  },
  {
    id: "cta-section",
    name: "CTA Section",
    description: "Call-to-action section with heading, text, and button",
    icon: Megaphone,
    category: "section",
    template: `## Ready to get started?

Join thousands of users who are already building amazing things with our platform.

<Button href="/sign-up">

Sign Up Now

</Button>
`,
  },
  {
    id: "comparison-columns",
    name: "Comparison Columns",
    description: "Side-by-side comparison using a two-column layout",
    icon: GitCompareArrows,
    category: "section",
    template: `<Columns cols={2}>

<Column>

### Option A

- Feature 1
- Feature 2
- Feature 3

</Column>

<Column>

### Option B

- Feature 1
- Feature 2
- Feature 3

</Column>

</Columns>
`,
  },
];

// ---------------------------------------------------------------------------
// Registry functions
// ---------------------------------------------------------------------------

export function getAllBlocks(): BlockDefinition[] {
  return blockDefinitions;
}

export function getBlocksByCategory(
  category: BlockCategory,
): BlockDefinition[] {
  return blockDefinitions.filter((block) => block.category === category);
}

export function searchBlocks(query: string): BlockDefinition[] {
  if (query === null || query === undefined || query === "") {
    return blockDefinitions;
  }

  const lowerQuery = query.toLowerCase();

  return blockDefinitions.filter((block) => {
    if (block.id.toLowerCase().includes(lowerQuery)) {
      return true;
    }
    if (block.name.toLowerCase().includes(lowerQuery)) {
      return true;
    }
    if (block.description.toLowerCase().includes(lowerQuery)) {
      return true;
    }
    return block.keywords.some((keyword) => keyword.toLowerCase().includes(lowerQuery));
  });
}

export function getAllPatterns(): BlockPattern[] {
  return blockPatterns;
}

export function searchPatterns(query: string): BlockPattern[] {
  if (query === null || query === undefined || query === "") {
    return blockPatterns;
  }

  const lowerQuery = query.toLowerCase();

  return blockPatterns.filter((pattern) => {
    if (pattern.id.toLowerCase().includes(lowerQuery)) {
      return true;
    }
    if (pattern.name.toLowerCase().includes(lowerQuery)) {
      return true;
    }
    return pattern.description.toLowerCase().includes(lowerQuery);
  });
}
