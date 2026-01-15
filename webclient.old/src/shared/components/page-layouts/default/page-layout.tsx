import * as React from "react";

import { ThemeProvider } from "@/shared/modules/themes/theme-provider.tsx";
import { backend } from "@/shared/modules/backend/backend.ts";
import type { Profile } from "@/shared/modules/backend/profiles/types.ts";

import { Header } from "./header.tsx";
import { Footer } from "./footer.tsx";
import { FontProvider } from "./font-provider.tsx";
import { ResponsiveIndicator } from "./responsive-indicator.tsx";
import "./page-layout.css";

type PageLayoutProps = {
  placeholders: Record<string, string>;
  children: React.ReactNode;
};

export async function PageLayout(props: PageLayoutProps) {
  let spotlight: Profile[] = [];
  try {
    spotlight = await backend.getSpotlight() ?? [];
  } catch (error) {
    // Silently handle backend errors - spotlight is optional
    console.error("Failed to fetch spotlight:", error);
  }

  return (
    <>
      <FontProvider />
      <ThemeProvider
        storageKey="theme"
        defaultTheme="default"
        enableColorScheme={false}
        disableTransitionOnChange
        themes={["default", "light", "midnight"]}
        attribute="data-theme"
      >
        <div className="flex flex-col min-h-screen">
          <Header placeholders={props.placeholders} spotlight={spotlight} />
          <main className="flex-1">{props.children}</main>
          <Footer />
          <ResponsiveIndicator />
        </div>
      </ThemeProvider>
    </>
  );
}
