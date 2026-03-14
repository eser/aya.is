import type { BlockCategory, BlockDefinition, BlockPattern } from "./types.ts";
import { generateContainerMdx, generateSelfClosingMdx } from "./mdx-template.ts";
import {
  AlertTriangle,
  AlignCenter,
  ArrowDownUp,
  ArrowRightLeft,
  BookOpen,
  ChevronDown,
  Code,
  Columns2,
  Columns3,
  CreditCard,
  ExternalLink,
  FileText,
  Frame,
  GitCompareArrows,
  Grid3X3,
  Heading2,
  HelpCircle,
  Image,
  Images,
  LayoutGrid,
  LayoutTemplate,
  Link,
  Megaphone,
  MessageSquareQuote,
  Minus,
  MousePointerClick,
  Music,
  PanelLeft,
  PanelTop,
  Quote,
  Space,
  Square,
  Table2,
  TextQuote,
  Twitter,
  Video,
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
    preview: '<Callout variant="info">...</Callout>',
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
    preview: '<Details summary="...">...</Details>',
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
    preview: "## Heading text",
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
    preview: "> Quote text",
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
    preview: "```js ... ```",
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

  {
    id: "table",
    name: "Table",
    description: "Structured data table",
    preview: "| Col | Col | Col |",
    icon: Table2,
    category: "text",
    keywords: ["table", "data", "grid", "rows", "columns"],
    props: [],
    generateMdx: () => {
      return `\n| Column 1 | Column 2 | Column 3 |\n| --- | --- | --- |\n| Cell 1 | Cell 2 | Cell 3 |\n| Cell 4 | Cell 5 | Cell 6 |\n`;
    },
  },
  {
    id: "pullquote",
    name: "Pullquote",
    description: "Emphasized quote with decorative styling",
    preview: '<Pullquote citation="...">...</Pullquote>',
    icon: TextQuote,
    category: "text",
    keywords: ["pullquote", "quote", "emphasis", "highlight", "cite"],
    props: [
      {
        name: "citation",
        type: "string",
        label: "Citation",
        required: false,
        defaultValue: "",
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (typeof values.citation === "string" && values.citation !== "") {
        props.citation = values.citation;
      }
      return generateContainerMdx("Pullquote", props, "Your emphasized quote here.");
    },
  },
  {
    id: "verse",
    name: "Verse",
    description: "Poetry and verse with preserved formatting",
    preview: "<Verse>line 1\\nline 2</Verse>",
    icon: BookOpen,
    category: "text",
    keywords: ["verse", "poetry", "poem", "lyrics", "preformatted"],
    props: [],
    generateMdx: () => {
      return generateContainerMdx("Verse", {}, "Line one\nLine two\nLine three");
    },
  },

  // ---- media ----
  {
    id: "image",
    name: "Image",
    description: "Display an image with alt text",
    preview: "![Alt](url)",
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

  {
    id: "gallery",
    name: "Gallery",
    description: "Image gallery grid",
    preview: "<Gallery cols={3}>images...</Gallery>",
    icon: Images,
    category: "media",
    keywords: ["gallery", "images", "photos", "grid"],
    props: [
      {
        name: "cols",
        type: "select",
        label: "Columns",
        required: false,
        defaultValue: "3",
        options: [
          { value: "2", label: "2" },
          { value: "3", label: "3" },
          { value: "4", label: "4" },
        ],
      },
    ],
    generateMdx: (values) => {
      const cols = typeof values.cols === "string" ? Number(values.cols) : 3;
      return `\n<Gallery cols={${cols}}>\n\n![Image 1](url1)\n![Image 2](url2)\n![Image 3](url3)\n\n</Gallery>\n`;
    },
  },
  {
    id: "audio",
    name: "Audio",
    description: "Audio player",
    preview: '<Audio src="..." />',
    icon: Music,
    category: "media",
    keywords: ["audio", "music", "sound", "podcast", "mp3"],
    props: [
      {
        name: "src",
        type: "string",
        label: "Source URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com/audio.mp3",
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
      if (typeof values.src === "string") {
        props.src = values.src;
      }
      if (typeof values.title === "string" && values.title !== "") {
        props.title = values.title;
      }
      return generateSelfClosingMdx("Audio", props);
    },
  },
  {
    id: "video",
    name: "Video",
    description: "Video player",
    preview: '<Video src="..." />',
    icon: Video,
    category: "media",
    keywords: ["video", "movie", "film", "mp4", "media"],
    props: [
      {
        name: "src",
        type: "string",
        label: "Source URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com/video.mp4",
      },
      {
        name: "poster",
        type: "string",
        label: "Poster URL",
        required: false,
        defaultValue: "",
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (typeof values.src === "string") {
        props.src = values.src;
      }
      if (typeof values.poster === "string" && values.poster !== "") {
        props.poster = values.poster;
      }
      return generateSelfClosingMdx("Video", props);
    },
  },
  {
    id: "cover",
    name: "Cover",
    description: "Background image with text overlay",
    preview: '<Cover src="..." overlay="dark">...</Cover>',
    icon: Frame,
    category: "media",
    keywords: ["cover", "hero", "banner", "background", "overlay"],
    props: [
      {
        name: "src",
        type: "string",
        label: "Source URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com/bg.jpg",
      },
      {
        name: "overlay",
        type: "select",
        label: "Overlay",
        required: false,
        defaultValue: "dark",
        options: [
          { value: "dark", label: "Dark" },
          { value: "light", label: "Light" },
          { value: "none", label: "None" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (typeof values.src === "string") {
        props.src = values.src;
      }
      if (values.overlay !== undefined && values.overlay !== "dark") {
        props.overlay = values.overlay;
      }
      return generateContainerMdx("Cover", props, "## Your heading here\n\nOverlay text content");
    },
  },

  // ---- layout ----
  {
    id: "columns-2",
    name: "2 Columns",
    description: "Two-column layout",
    preview: "<Columns cols={2}><Column>...</Column></Columns>",
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
    preview: "<Columns cols={3}>...</Columns>",
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
    preview: "<Divider />",
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
    preview: '<Spacer size="md" />',
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
    preview: "<Cards><Card />...</Cards>",
    icon: LayoutGrid,
    category: "layout",
    keywords: ["cards", "grid", "card grid", "card list"],
    props: [],
    generateMdx: () => {
      return `\n<Cards>\n\n<Card title="Card 1" description="Description" />\n\n<Card title="Card 2" description="Description" />\n\n</Cards>\n`;
    },
  },

  {
    id: "media-text",
    name: "Media & Text",
    description: "Image alongside text content",
    preview: '<MediaText src="...">text</MediaText>',
    icon: PanelLeft,
    category: "layout",
    keywords: ["media", "text", "image", "side", "split", "layout"],
    props: [
      {
        name: "src",
        type: "string",
        label: "Source URL",
        required: true,
        defaultValue: "",
        placeholder: "https://example.com/image.jpg",
      },
      {
        name: "mediaPosition",
        type: "select",
        label: "Media Position",
        required: false,
        defaultValue: "left",
        options: [
          { value: "left", label: "Left" },
          { value: "right", label: "Right" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (typeof values.src === "string") {
        props.src = values.src;
      }
      if (values.mediaPosition !== undefined && values.mediaPosition !== "left") {
        props.mediaPosition = values.mediaPosition;
      }
      return generateContainerMdx("MediaText", props, "Your text content alongside the media.");
    },
  },
  {
    id: "row",
    name: "Row",
    description: "Horizontal flex layout",
    preview: '<Row gap="md">...</Row>',
    icon: ArrowRightLeft,
    category: "layout",
    keywords: ["row", "horizontal", "flex", "inline", "layout"],
    props: [
      {
        name: "gap",
        type: "select",
        label: "Gap",
        required: false,
        defaultValue: "md",
        options: [
          { value: "sm", label: "Small" },
          { value: "md", label: "Medium" },
          { value: "lg", label: "Large" },
        ],
      },
      {
        name: "justify",
        type: "select",
        label: "Justify",
        required: false,
        defaultValue: "start",
        options: [
          { value: "start", label: "Start" },
          { value: "center", label: "Center" },
          { value: "end", label: "End" },
          { value: "between", label: "Between" },
          { value: "around", label: "Around" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.gap !== undefined && values.gap !== "md") {
        props.gap = values.gap;
      }
      if (values.justify !== undefined && values.justify !== "start") {
        props.justify = values.justify;
      }
      return generateContainerMdx("Row", props, "Content items here...");
    },
  },
  {
    id: "stack",
    name: "Stack",
    description: "Vertical flex layout",
    preview: '<Stack gap="md">...</Stack>',
    icon: ArrowDownUp,
    category: "layout",
    keywords: ["stack", "vertical", "flex", "layout"],
    props: [
      {
        name: "gap",
        type: "select",
        label: "Gap",
        required: false,
        defaultValue: "md",
        options: [
          { value: "sm", label: "Small" },
          { value: "md", label: "Medium" },
          { value: "lg", label: "Large" },
        ],
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.gap !== undefined && values.gap !== "md") {
        props.gap = values.gap;
      }
      return generateContainerMdx("Stack", props, "Stacked content here...");
    },
  },
  {
    id: "grid",
    name: "Grid",
    description: "Responsive grid layout",
    preview: "<Grid cols={2}>...</Grid>",
    icon: Grid3X3,
    category: "layout",
    keywords: ["grid", "auto", "responsive", "layout", "tiles"],
    props: [
      {
        name: "cols",
        type: "select",
        label: "Columns",
        required: false,
        defaultValue: "2",
        options: [
          { value: "2", label: "2" },
          { value: "3", label: "3" },
          { value: "4", label: "4" },
        ],
      },
      {
        name: "gap",
        type: "select",
        label: "Gap",
        required: false,
        defaultValue: "md",
        options: [
          { value: "sm", label: "Small" },
          { value: "md", label: "Medium" },
          { value: "lg", label: "Large" },
        ],
      },
    ],
    generateMdx: (values) => {
      const cols = typeof values.cols === "string" ? Number(values.cols) : 2;
      const props: Record<string, string | number | boolean> = { cols };
      if (values.gap !== undefined && values.gap !== "md") {
        props.gap = values.gap;
      }
      return generateContainerMdx("Grid", props, "Grid items here...");
    },
  },
  {
    id: "group",
    name: "Group",
    description: "Container with visual styling",
    preview: '<Group padding="md">...</Group>',
    icon: Square,
    category: "layout",
    keywords: ["group", "container", "box", "section", "wrapper"],
    props: [
      {
        name: "padding",
        type: "select",
        label: "Padding",
        required: false,
        defaultValue: "md",
        options: [
          { value: "none", label: "None" },
          { value: "sm", label: "Small" },
          { value: "md", label: "Medium" },
          { value: "lg", label: "Large" },
        ],
      },
      {
        name: "background",
        type: "select",
        label: "Background",
        required: false,
        defaultValue: "none",
        options: [
          { value: "none", label: "None" },
          { value: "muted", label: "Muted" },
          { value: "card", label: "Card" },
          { value: "accent", label: "Accent" },
        ],
      },
      {
        name: "border",
        type: "boolean",
        label: "Border",
        required: false,
        defaultValue: false,
      },
      {
        name: "rounded",
        type: "boolean",
        label: "Rounded",
        required: false,
        defaultValue: true,
      },
    ],
    generateMdx: (values) => {
      const props: Record<string, string | number | boolean> = {};
      if (values.padding !== undefined && values.padding !== "md") {
        props.padding = values.padding;
      }
      if (values.background !== undefined && values.background !== "none") {
        props.background = values.background;
      }
      if (values.border === true) {
        props.border = true;
      }
      if (values.rounded === false) {
        props.rounded = false;
      }
      return generateContainerMdx("Group", props, "Grouped content here...");
    },
  },
  {
    id: "center",
    name: "Center",
    description: "Center-aligned content",
    preview: "<Center>...</Center>",
    icon: AlignCenter,
    category: "layout",
    keywords: ["center", "align", "middle", "layout"],
    props: [],
    generateMdx: () => {
      return generateContainerMdx("Center", {}, "Centered content here...");
    },
  },

  // ---- data ----
  {
    id: "card",
    name: "Card",
    description: "Content card with title, description, and optional link",
    preview: '<Card title="..." description="..." />',
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
    preview: '<SiteLink href="...">text</SiteLink>',
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
    preview: '<YouTubeEmbed videoId="..." />',
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
    preview: '<TwitterEmbed tweetId="..." />',
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
    preview: '<PDF src="..." />',
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
    preview: '<Embed url="..." />',
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
    preview: '<Button href="...">text</Button>',
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
    preview: '<Tabs><Tab label="...">...</Tab></Tabs>',
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
