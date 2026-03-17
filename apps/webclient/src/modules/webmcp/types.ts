// Copyright 2023-present Eser Ozvataf and other contributors. All rights reserved. Apache-2.0 license.
// WebMCP API type declarations
// Spec: https://webmachinelearning.github.io/webmcp/
// These types augment the Navigator interface since the API is not yet in lib.dom.d.ts

export interface ToolAnnotations {
  readOnlyHint?: boolean;
}

export interface ModelContextClient {
  requestUserInteraction: (
    callback: () => Promise<unknown>,
  ) => Promise<unknown>;
}

export interface ModelContextTool {
  name: string;
  description: string;
  inputSchema?: object;
  execute: (
    input: Record<string, unknown>,
    client: ModelContextClient,
  ) => Promise<unknown>;
  annotations?: ToolAnnotations;
}

export interface ModelContextOptions {
  tools?: ModelContextTool[];
}

export interface ModelContext {
  provideContext(options?: ModelContextOptions): void;
  clearContext(): void;
  registerTool(tool: ModelContextTool): void;
  unregisterTool(name: string): void;
}

declare global {
  interface Navigator {
    modelContext?: ModelContext;
  }
}
