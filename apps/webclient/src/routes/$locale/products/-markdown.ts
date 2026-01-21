/**
 * Products domain - markdown generation utilities
 * Products shows profiles with kind="product"
 */
import { formatFrontmatter } from "@/lib/markdown";
import { registerMarkdownHandler } from "@/lib/markdown-middleware";
import { backend } from "@/modules/backend/backend";
import type { Profile } from "@/modules/backend/types";

/**
 * Generate markdown for the products listing page
 */
export function generateProductsListingMarkdown(
  products: Profile[],
  locale: string,
): string {
  const frontmatter = formatFrontmatter({
    title: "Products",
    generated_at: new Date().toISOString(),
  });

  const productLinks = products.map((product) => {
    let line = `- [${product.title}](/${locale}/${product.slug}.md)`;
    if (product.description !== null && product.description !== undefined) {
      line += `\n  ${product.description.slice(0, 200)}${product.description.length > 200 ? "..." : ""}`;
    }
    return line;
  });

  return `${frontmatter}\n\n${productLinks.join("\n\n")}`;
}

/**
 * Register markdown handler for products listing
 * Pattern: /$locale/products
 */
export function registerProductsListingHandler(): void {
  registerMarkdownHandler("$locale/products", async (_params, locale, _searchParams) => {
    const products = await backend.getProfilesByKinds(locale, ["product"]);

    if (products === null) {
      return null;
    }

    return generateProductsListingMarkdown(products, locale);
  });
}
