// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
import { defineConfig } from "vitest/config";
import tsConfigPaths from "vite-tsconfig-paths";

// Kept separate from vite.config.ts on purpose: tests must not load the TanStack Start / nitro / devtools plugins.
export default defineConfig({
  plugins: [
    tsConfigPaths({
      projects: ["./tsconfig.json"],
    }),
  ],
  test: {
    include: ["src/**/*.test.ts"],
    environment: "node",
  },
});
