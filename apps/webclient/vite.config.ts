// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import { nitro } from "nitro/vite";
import tsConfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";

export default defineConfig({
  resolve: {
    alias: {
      // vite-tsconfig-paths only rewrites JS/TS specifiers; CSS and asset imports
      // need an explicit alias so `@/theme.css` resolves during the build.
      "@": fileURLToPath(new URL("./src", import.meta.url)),

      // speedy-router (https://github.com/anonrig/router) is a drop-in for the TanStack Router runtime.
      // Source imports stay `@tanstack/react-router`; remove these three lines, the `ssr` block below and the
      // speedy-router* entries of `optimizeDeps.include` together with the speedy-router* dependencies to switch back.
      "@tanstack/react-router": "speedy-router",
      "@tanstack/router-core": "speedy-router-core",
      "@tanstack/history": "speedy-router-history",
    },
  },
  ssr: {
    // Aliases only apply to code that Vite processes. Bundle every @tanstack/* package on the server too, otherwise
    // TanStack Start would load its own externalized copy of the real router next to speedy-router.
    noExternal: [/^@tanstack\//],
  },
  optimizeDeps: {
    // TanStack Start has no index.html for Vite's dependency scan to crawl, so these were discovered only once the browser
    // asked for them. That triggered a mid-session re-bundle ("optimized dependencies changed. reloading") that deleted the
    // chunks other requests were still loading. Listing them up front makes the first optimize pass complete.
    include: [
      "speedy-router",
      "speedy-router-core",
      "speedy-router-core/is-server",
      "speedy-router-history",
      "@tanstack/router-core/isServer",
      "@tanstack/router-core/ssr/client",
    ],
  },
  server: {
    port: 3000,
    host: "::",
    watch: {
      ignored: ["**/routeTree.gen.ts"],
    },
    allowedHosts: true,
  },
  plugins: [
    devtools(),
    nitro(),
    tsConfigPaths({
      projects: ["./tsconfig.json"],
    }),
    tailwindcss(),
    tanstackStart({
      srcDirectory: "src",
      router: {
        routeFileIgnorePattern: "^_components/",
        routeTreeFileHeader: [
          "// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.",
          "/* eslint-disable */",
          "// @ts-nocheck",
          "// noinspection JSUnusedGlobalSymbols",
        ],
      },
    }),
    viteReact(),
  ],
});
