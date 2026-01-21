import process from "node:process";
import { createMiddleware, createStart } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
import type { RequestContext } from "@/request-context";
import { getDomainConfiguration } from "./server/domain-configurations";
import { requestContextBinder } from "./server/request-context-binder";
import { markdownMiddleware } from "./lib/markdown-middleware";
import { registerAllMarkdownHandlers } from "./lib/markdown-handlers";

// Register all markdown handlers at startup
registerAllMarkdownHandlers();

const customDomainMiddleware = createMiddleware()
  .server(async ({ request, next }) => {
    const url = new URL(request.url);

    // Skip static assets and internal paths - pass through directly
    if (
      url.pathname.startsWith("/_") ||
      url.pathname.startsWith("/api/") ||
      url.pathname.startsWith("/assets/") ||
      url.pathname.includes(".")
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
      pathParts.splice(1, 0, domainConfiguration.profileSlug);
    }

    const requestContext: RequestContext = {
      domainConfiguration: domainConfiguration,
      path: pathParts,
      originalPath: originalPathParts,
    };

    return requestContextBinder.run(requestContext, () => {
      return next({ context: { requestContext } });
    });
  });

export const startInstance = createStart(() => {
  return {
    // Markdown middleware runs first to intercept .md requests
    requestMiddleware: [markdownMiddleware, customDomainMiddleware],
  };
});
