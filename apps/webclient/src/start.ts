import process from "node:process";
import { createMiddleware, createStart } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
import type { RequestContext } from "@/request-context";
import { getDomainConfiguration } from "./server/domain-configurations";
import { markdownMiddleware } from "./server/markdown-middleware";
import { registerAllMarkdownHandlers } from "./server/markdown-handlers";

// Register all markdown handlers at startup
registerAllMarkdownHandlers();

const customDomainMiddleware = createMiddleware()
  .server(async ({ request, next }) => {
    const url = new URL(request.url);

    // Skip static assets and internal paths - pass through directly
    // Note: .md files are NOT skipped - they need path rewriting for custom domains
    if (
      url.pathname.startsWith("/_") ||
      url.pathname.startsWith("/api/") ||
      url.pathname.startsWith("/assets/") ||
      (url.pathname.includes(".") && !url.pathname.endsWith(".md") && !url.pathname.endsWith(".txt"))
    ) {
      return next();
    }

    const forwardedHost = getRequestHeader("x-forwarded-host")?.split(":")[0];
    const hostHeader = getRequestHeader("host")?.split(":")[0];
    const hostFromHeaders = forwardedHost ?? hostHeader ?? url.hostname;
    const hostname = process.env.CUSTOM_DOMAIN ?? hostFromHeaders;

    // console.log("forwardedHost", forwardedHost);
    // console.log("hostHeader", hostHeader);
    // console.log("url.hostname", url.hostname);
    // console.log("hostname", hostname);

    // Get domain configuration (static or from backend)
    const domainConfiguration = await getDomainConfiguration(hostname);

    const originalPathParts = url.pathname.split("/").filter(Boolean);
    const pathParts = [...originalPathParts];
    if (domainConfiguration.type === "custom-domain") {
      // Handle the case where the path is just /{locale}.md (e.g., /tr.md)
      // In this case, we need to transform it to /{locale}/{profileSlug}.md
      if (pathParts.length === 1 && pathParts[0] !== undefined && pathParts[0].endsWith(".md")) {
        const localeWithMd = pathParts[0];
        const locale = localeWithMd.slice(0, -3); // Remove .md suffix
        pathParts[0] = locale;
        pathParts.push(`${domainConfiguration.profileSlug}.md`);
      } else {
        // Standard case: insert profile slug after locale
        pathParts.splice(1, 0, domainConfiguration.profileSlug);
      }
    }

    const requestContext: RequestContext = {
      domainConfiguration: domainConfiguration,
      path: pathParts,
      originalPath: originalPathParts,
    };

    const { requestContextBinder } = await import("./server/request-context-binder");
    return requestContextBinder.run(requestContext, () => {
      return next({ context: { requestContext } });
    });
  });

export const startInstance = createStart(() => {
  return {
    // Custom domain middleware runs first to set up request context with rewritten paths
    // Markdown middleware then uses the rewritten path for custom domains
    requestMiddleware: [customDomainMiddleware, markdownMiddleware],
  };
});
