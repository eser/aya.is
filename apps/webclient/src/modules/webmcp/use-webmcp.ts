import { useEffect } from "react";
import { useNavigate, useRouterState } from "@tanstack/react-router";
import "./types.ts"; // Side-effect import: augments Navigator with modelContext
import { buildContextualTools, buildGlobalTools, parseRouteContext } from "./tools.ts";

export function useWebMCP(locale: string, isCustomDomain: boolean): void {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const navigate = useNavigate();

  useEffect(() => {
    // Feature detection: no-op if WebMCP API is unavailable
    if (navigator.modelContext === undefined) {
      return;
    }

    const context = parseRouteContext(pathname);
    const globalTools = buildGlobalTools(locale, isCustomDomain, navigate);
    const contextualTools = buildContextualTools(locale, context);

    // provideContext clears previous tools and registers the new set atomically,
    // which is the recommended pattern for SPAs that change tools on navigation.
    navigator.modelContext.provideContext({
      tools: [...globalTools, ...contextualTools],
    });

    return () => {
      navigator.modelContext?.clearContext();
    };
  }, [locale, isCustomDomain, pathname, navigate]);
}
