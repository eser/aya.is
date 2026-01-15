import { defineConfig } from "vite";
import { devtools } from "@tanstack/devtools-vite";
import { nitro } from "nitro/vite";
import tsConfigPaths from "vite-tsconfig-paths";
import tailwindcss from "@tailwindcss/vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";

export default defineConfig({
  server: {
    port: 3000,
    host: "::",
    watch: {
      ignored: ["**/routeTree.gen.ts"],
    },
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
      tsr: {
        routeFileIgnorePattern: "^_components/",
      },
    }),
    viteReact(),
  ],
});
