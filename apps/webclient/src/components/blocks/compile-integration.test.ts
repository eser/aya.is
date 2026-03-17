// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
/// <reference lib="deno.ns" />
import { compile } from "@mdx-js/mdx";
import { getAllBlocks, getAllPatterns } from "./registry.ts";

/**
 * Lightweight MDX compile helper for tests.
 * Uses the same outputFormat as compileMdxLite but without heavy plugins
 * that pull in CSS-module-importing React components.
 */
async function compileTestMdx(source: string): Promise<string> {
  const compiled = await compile(source, {
    outputFormat: "function-body",
  });
  return String(compiled);
}

Deno.test("every block template compiles", async () => {
  const blocks = getAllBlocks();
  for (const block of blocks) {
    const defaultValues: Record<string, string | number | boolean> = {};
    for (const prop of block.props) {
      defaultValues[prop.name] = prop.defaultValue;
    }
    const mdx = block.generateMdx(defaultValues);
    await compileTestMdx(mdx);
  }
});

Deno.test("every pattern template compiles", async () => {
  const patterns = getAllPatterns();
  for (const pattern of patterns) {
    await compileTestMdx(pattern.template);
  }
});

Deno.test("all blocks concatenated compile", async () => {
  const blocks = getAllBlocks();
  const allMdx = blocks.map((block) => {
    const defaultValues: Record<string, string | number | boolean> = {};
    for (const prop of block.props) {
      defaultValues[prop.name] = prop.defaultValue;
    }
    return block.generateMdx(defaultValues);
  }).join("\n\n");
  await compileTestMdx(allMdx);
});
