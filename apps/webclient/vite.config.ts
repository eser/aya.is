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
    },
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
